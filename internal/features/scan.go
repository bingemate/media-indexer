package features

import (
	"errors"
	"fmt"
	objectStorage "github.com/bingemate/media-go-pkg/object-storage"
	"github.com/bingemate/media-indexer/internal/repository"
	"github.com/bingemate/media-indexer/pkg"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

// MovieScanner represents a struct that scans movie folders to search for movie files and move them.
type MovieScanner struct {
	source          string                      // Source directory path to scan for movies.
	destination     string                      // Destination directory path to move the found movies.
	mediaClient     pkg.MediaClient             // Media client object to search for movies on TMDB.
	mediaRepository *repository.MediaRepository // Media repository object to save the media files and their details.
	objectStorage   objectStorage.ObjectStorage // Object storage object to upload the media files.
}

// MovieScannerResult represents a struct that holds the results of scanning and moving movie files.
type MovieScannerResult struct {
	Source string    // Source filename.
	Movie  pkg.Movie // Movie details returned by TMDB.
}

// TVScanner represents a struct that scans TV show folders to search for TV show files and move them.
type TVScanner struct {
	source          string                      // Source directory path to scan for TV shows.
	destination     string                      // Destination directory path to move the found TV shows.
	mediaClient     pkg.MediaClient             // Media client object to search for TV shows on TMDB.
	mediaRepository *repository.MediaRepository // Media repository object to save the media files and their details.
	objectStorage   objectStorage.ObjectStorage // Object storage object to upload the media files.
}

// TVScannerResult represents a struct that holds the results of scanning and moving TV show files.
type TVScannerResult struct {
	Source    string        // Source filename.
	TVEpisode pkg.TVEpisode // TV episode details returned by TMDB.
}

// NewMovieScanner returns a new instance of MovieScanner with given source directory, target directory, and TMDB API key.
func NewMovieScanner(source, destination string, mediaClient pkg.MediaClient, mediaRepository *repository.MediaRepository, objectStorage objectStorage.ObjectStorage) *MovieScanner {
	return &MovieScanner{
		source:          source,
		destination:     destination,
		mediaClient:     mediaClient,
		mediaRepository: mediaRepository,
		objectStorage:   objectStorage,
	}
}

// NewTVScanner returns a new instance of TVScanner with given source directory, target directory, and TMDB API key.
func NewTVScanner(source, destination string, mediaClient pkg.MediaClient, mediaRepository *repository.MediaRepository, objectStorage objectStorage.ObjectStorage) *TVScanner {
	return &TVScanner{
		source:          source,
		destination:     destination,
		mediaClient:     mediaClient,
		mediaRepository: mediaRepository,
		objectStorage:   objectStorage,
	}
}

// ScanMovies scans the source directory for movies and moves them to the destination directory.
// It returns a slice of MovieScannerResult and an error if there is any.
func (s *MovieScanner) ScanMovies() error {
	// Locks the scanner to prevent concurrent scanning
	locked := jobLock.TryLock()
	if !locked {
		log.Printf("Job '%s' already running, skipping this run", pkg.GetJobName())
		return fmt.Errorf("job '%s' already running, skipping this run", pkg.GetJobName())
	}

	go func() {
		defer jobLock.Unlock()
		pkg.ClearJobLogs("scan movies")

		mediaFiles, err := s.scanMovieFolder()
		if err != nil {
			log.Printf("Failed to scan movie folder: %v", err)
			pkg.AppendJobLog(fmt.Sprintf("Failed to scan movie folder: %v", err))
			return
		}
		atomicMovieList := s.retrieveMovieList(mediaFiles)

		result := s.buildMovieScannerResult(atomicMovieList)

		// Process the movies to the destination directory and returns an error if it fails
		err = s.processMovies(atomicMovieList, s.destination)
		if err != nil {
			log.Printf("Failed to process movies to %s: %v", s.destination, err)
			pkg.AppendJobLog(fmt.Sprintf("Failed to process movies to %s: %v", s.destination, err))
		}

		log.Printf("Processed %d movies to %s.", len(*result), s.destination)
		pkg.AppendJobLog(fmt.Sprintf("Processed %d movies to %s.", len(*result), s.destination))

		/*err = pkg.ClearFolderContent(s.source)
		if err != nil {
			log.Printf("Failed to clear source folder content: %v", err)
			pkg.AppendJobLog(fmt.Sprintf("Failed to clear source folder content: %v", err))
			return
		}*/
	}()
	return nil
}

func (s *MovieScanner) buildMovieScannerResult(atomicMediaList *pkg.AtomicMovieList) *[]MovieScannerResult {
	// Initializes an empty slice of MovieScannerResult
	var result = make([]MovieScannerResult, 0)

	// Iterates through each media file and its corresponding movie information in the AtomicMovieList and adds it to the result slice
	for mediaFile, media := range atomicMediaList.GetAll() {
		result = append(result, MovieScannerResult{
			Source: mediaFile.Filename,
			Movie:  media,
		})
	}
	return &result
}

func (s *MovieScanner) scanMovieFolder() (*[]pkg.MovieFile, error) {
	// Logs that the function is scanning the source directory for movies
	log.Printf("Scanning %s for movies...", s.source)
	pkg.AppendJobLog(fmt.Sprintf("Scanning %s for movies...", s.source))

	// Builds the directory tree from the source directory and returns an error if it fails
	mediaFiles, err := pkg.BuildMovieTree(s.source)
	if err != nil {
		log.Printf("Failed to scan source tree: %v", err)
		pkg.AppendJobLog(fmt.Sprintf("Failed to scan source tree: %v", err))
		return nil, err
	}

	log.Printf("Scanning %d files in %s...", len(mediaFiles), s.source)
	pkg.AppendJobLog(fmt.Sprintf("Scanning %d files in %s...", len(mediaFiles), s.source))
	return &mediaFiles, nil
}

func (s *MovieScanner) retrieveMovieList(mediaFiles *[]pkg.MovieFile) *pkg.AtomicMovieList {
	// Initialize a WaitGroup and an AtomicMovieList
	var wg sync.WaitGroup
	var atomicMovieList = pkg.NewAtomicMovieList()

	// Create a semaphore channel to limit the number of goroutines
	sem := make(chan bool, 4)

	for _, mediaFile := range *mediaFiles {
		sem <- true
		wg.Add(1)
		go func(mediaFile pkg.MovieFile) {
			defer wg.Done()
			defer func() { <-sem }()

			log.Printf("Searching for movie information for file %s...", mediaFile.Filename)
			pkg.AppendJobLog(fmt.Sprintf("Searching for movie information for file %s...", mediaFile.Filename))

			media, ok := searchMovie(&mediaFile, s.mediaClient)
			if !ok {
				log.Printf("Failed to find movie information for file %s.", mediaFile.Filename)
				pkg.AppendJobLog(fmt.Sprintf("Failed to find movie information for file %s.", mediaFile.Filename))
				return
			}
			atomicMovieList.LinkMediaFile(mediaFile, media)
		}(mediaFile)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	wg.Wait()

	log.Println("Movie scan complete.")
	pkg.AppendJobLog("Movie scan complete.")
	return atomicMovieList
}

// ScanTV scans the source directory for TV shows and moves them to the destination directory.
// It returns a slice of TVScannerResult and an error if there is any.
func (s *TVScanner) ScanTV() error {
	// Locks the scanner to prevent concurrent scanning
	locked := jobLock.TryLock()
	if !locked {
		log.Printf("Job '%s' already running, skipping this run", pkg.GetJobName())
		return fmt.Errorf("job '%s' already running, skipping this run", pkg.GetJobName())
	}
	pkg.ClearJobLogs("scan tv")

	go func() {

		defer jobLock.Unlock()

		mediaFiles, err := s.scanTVFolder()
		if err != nil {
			log.Printf("Failed to scan TV folder: %v", err)
			pkg.AppendJobLog(fmt.Sprintf("Failed to scan TV folder: %v", err))
			return
		}

		atomicMediaList := s.retrieveTvList(mediaFiles)

		result := s.buildTVScannerResult(atomicMediaList)

		// Moves the TV shows to the destination directory and returns an error if it fails
		err = s.processTVEpisodes(atomicMediaList, s.destination)
		if err != nil {
			log.Printf("Failed to process TV shows to %s: %v", s.destination, err)
			pkg.AppendJobLog(fmt.Sprintf("Failed to process TV shows to %s: %v", s.destination, err))
		}

		log.Printf("Processed %d TV shows to %s.", len(*result), s.destination)
		pkg.AppendJobLog(fmt.Sprintf("Processed %d TV shows to %s.", len(*result), s.destination))

		/*err = pkg.ClearFolderContent(s.source)
		if err != nil {
			log.Printf("Failed to clear source folder content: %v", err)
			pkg.AppendJobLog(fmt.Sprintf("Failed to clear source folder content: %v", err))
			return
		}*/
	}()
	return nil
}

func (s *TVScanner) buildTVScannerResult(atomicMediaList *pkg.AtomicTVEpisodeList) *[]TVScannerResult {
	// Initializes an empty slice of TVScannerResult
	var result = make([]TVScannerResult, 0)

	// Iterates through each media file and its corresponding TV show information in the AtomicMovieList and adds it to the result slice
	for mediaFile, media := range atomicMediaList.GetAll() {
		result = append(result, TVScannerResult{
			Source:    mediaFile.Filename,
			TVEpisode: media,
		})
	}
	return &result
}

func (s *TVScanner) retrieveTvList(mediaFiles *[]pkg.TVShowFile) *pkg.AtomicTVEpisodeList {
	var wg sync.WaitGroup
	var atomicMediaList = pkg.NewAtomicTVEpisodeList()

	// Create a semaphore channel to limit the number of goroutines
	sem := make(chan bool, 4)

	for _, mediaFile := range *mediaFiles {
		sem <- true
		wg.Add(1)
		go func(mediaFile pkg.TVShowFile) {
			defer wg.Done()
			defer func() { <-sem }()

			log.Printf("Searching for TV show information for file %s...", mediaFile.Filename)
			pkg.AppendJobLog(fmt.Sprintf("Searching for TV show information for file %s...", mediaFile.Filename))

			media, ok := searchTVEpisode(&mediaFile, s.mediaClient)
			if !ok {
				log.Printf("Failed to find TV show information for file %s.", mediaFile.Filename)
				pkg.AppendJobLog(fmt.Sprintf("Failed to find TV show information for file %s.", mediaFile.Filename))
				return
			}
			log.Printf("Found TV show information for file %s:", mediaFile.Filename)
			pkg.AppendJobLog(fmt.Sprintf("Found TV show information for file %s:", mediaFile.Filename))
			log.Println(media)
			pkg.AppendJobLog(fmt.Sprintf("%v", media))
			atomicMediaList.LinkMediaFile(mediaFile, media)
		}(mediaFile)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	wg.Wait()

	log.Println("TV show scan complete.")
	pkg.AppendJobLog("TV show scan complete.")
	return atomicMediaList
}

func (s *TVScanner) scanTVFolder() (*[]pkg.TVShowFile, error) {
	// Logs that the function is scanning the source directory for TV shows
	log.Printf("Scanning %s for TV shows...", s.source)
	pkg.AppendJobLog(fmt.Sprintf("Scanning %s for TV shows...", s.source))

	// Builds the directory tree from the source directory and returns an error if it fails
	mediaFiles, err := pkg.BuildTVShowTree(s.source)
	if err != nil {
		log.Printf("Failed to scan source tree: %v", err)
		pkg.AppendJobLog(fmt.Sprintf("Failed to scan source tree: %v", err))
		return nil, err
	}

	log.Printf("Scanning %d files in %s...", len(mediaFiles), s.source)
	pkg.AppendJobLog(fmt.Sprintf("Scanning %d files in %s...", len(mediaFiles), s.source))
	return &mediaFiles, nil
}

// searchMovie searches for a movie on TMDB using the media file name and year, returning the movie details and a boolean indicating whether it was found.
func searchMovie(mediaFile *pkg.MovieFile, client pkg.MediaClient) (pkg.Movie, bool) {
	result, err := client.SearchMovie(mediaFile.SanitizedName, mediaFile.Year)
	if err != nil {
		log.Printf("Error while media search on %s : %s. Sanitized name was : %s", mediaFile.Filename, err.Error(), mediaFile.SanitizedName)
		pkg.AppendJobLog(fmt.Sprintf("Error while media search on %s : %s. Sanitized name was : %s", mediaFile.Filename, err.Error(), mediaFile.SanitizedName))
		return pkg.Movie{}, false
	}
	return result, true
}

// searchTVEpisode searches for a TV show on TMDB using the media file name and year, returning the TV show details and a boolean indicating whether it was found.
func searchTVEpisode(mediaFile *pkg.TVShowFile, client pkg.MediaClient) (pkg.TVEpisode, bool) {
	result, err := client.SearchTVShow(mediaFile.SanitizedName, mediaFile.Season, mediaFile.Episode)
	if err != nil {
		log.Printf("Error while media search on %s : %s. Sanitized name was : %s", mediaFile.Filename, err.Error(), mediaFile.SanitizedName)
		pkg.AppendJobLog(fmt.Sprintf("Error while media search on %s : %s. Sanitized name was : %s", mediaFile.Filename, err.Error(), mediaFile.SanitizedName))
		return pkg.TVEpisode{}, false
	}
	return result, true
}

// processMovies moves the media files to the destination directory path provided as argument.
// It returns an error if the destination directory does not exist or if there was an error while moving the file.
func (s *MovieScanner) processMovies(movieList *pkg.AtomicMovieList, destination string) error {
	if !pkg.IsDirectoryExists(destination) {
		pkg.AppendJobLog(fmt.Sprintf("Destination directory %s does not exists", destination))
		return errors.New("destination directory does not exists")
	}
	var now time.Time
	for mediaFile, media := range movieList.GetAll() {
		now = time.Now()
		var source = path.Join(mediaFile.Path, mediaFile.Filename)
		err := s.mediaRepository.IndexMovie(media, source, s.destination)
		if err != nil {
			log.Printf("Failed to index %s to %s : %s", source, s.destination, err.Error())
			pkg.AppendJobLog(fmt.Sprintf("Failed to index %s to %s : %s", source, s.destination, err.Error()))
			return err
		}
		log.Printf("Processed %s - %s %s. Took %v", mediaFile.Filename, media.Name, media.Year(), time.Since(now))
		pkg.AppendJobLog(fmt.Sprintf("Processed %-60s - %s %s. Took %v", mediaFile.Filename, media.Name, media.Year(), time.Since(now)))

		go func(mediaFile pkg.MovieFile, media pkg.Movie, destination string) {
			log.Printf("Removing %s", source)
			pkg.AppendJobLog(fmt.Sprintf("Removing %s", source))
			err = os.Remove(source)
			if err != nil {
				log.Printf("Failed to remove %s : %s", source, err.Error())
				pkg.AppendJobLog(fmt.Sprintf("Failed to remove %s : %s", source, err.Error()))
			}
			// Upload destination to S3
			now = time.Now()
			log.Printf("Uploading movie %d to S3...", media.ID)
			pkg.AppendJobLog(fmt.Sprintf("Uploading movie %d to S3...", media.ID))
			err = s.objectStorage.UploadMediaFiles(
				path.Join("movies", strconv.Itoa(media.ID)),
				path.Join(destination, strconv.Itoa(media.ID)),
			)
			if err != nil {
				log.Printf("Failed to upload %s to S3 : %s", destination, err.Error())
				pkg.AppendJobLog(fmt.Sprintf("Failed to upload %s to S3 : %s", destination, err.Error()))
			} else {
				log.Printf("Uploaded %s to S3. Took %v", destination, time.Since(now))
				pkg.AppendJobLog(fmt.Sprintf("Uploaded %s to S3. Took %v", destination, time.Since(now)))
			}
			// Remove destination from local
			log.Printf("Removing %s from local storage", path.Join(destination, strconv.Itoa(media.ID)))
			pkg.AppendJobLog(fmt.Sprintf("Removing %s from local storage", path.Join(destination, strconv.Itoa(media.ID))))
			err = os.RemoveAll(path.Join(destination, strconv.Itoa(media.ID)))
			if err != nil {
				log.Printf("Failed to remove %s : %s", path.Join(destination, strconv.Itoa(media.ID)), err.Error())
				pkg.AppendJobLog(fmt.Sprintf("Failed to remove %s : %s", path.Join(destination, strconv.Itoa(media.ID)), err.Error()))
			} else {
				log.Printf("Removed %s from local storage", path.Join(destination, strconv.Itoa(media.ID)))
				pkg.AppendJobLog(fmt.Sprintf("Removed %s from local storage", path.Join(destination, strconv.Itoa(media.ID))))
			}
		}(mediaFile, media, destination)
	}
	return nil
}

// processTVEpisodes moves the media files to the destination directory path provided as argument.
// It returns an error if the destination directory does not exist or if there was an error while moving the file.
func (s *TVScanner) processTVEpisodes(tvList *pkg.AtomicTVEpisodeList, destination string) error {
	if !pkg.IsDirectoryExists(destination) {
		pkg.AppendJobLog(fmt.Sprintf("Destination directory %s does not exists", destination))
		return errors.New("destination directory does not exists")
	}
	var now time.Time
	for mediaFile, media := range tvList.GetAll() {
		now = time.Now()
		var source = path.Join(mediaFile.Path, mediaFile.Filename)
		err := s.mediaRepository.IndexTvEpisode(media, source, s.destination)
		if err != nil {
			log.Printf("Failed to index %s to %s : %s", source, s.destination, err.Error())
			pkg.AppendJobLog(fmt.Sprintf("Failed to index %s to %s : %s", source, s.destination, err.Error()))
			return err
		}
		log.Printf("Processed %-60s - %s - %s s%02de%02d\nTook %s", mediaFile.Filename, media.TvShowName, media.Year(), mediaFile.Season, mediaFile.Episode, time.Since(now))
		pkg.AppendJobLog(fmt.Sprintf("Processed %-60s - %s - %s\nTook %s", mediaFile.Filename, media.TvShowName, media.Year(), time.Since(now)))

		go func(mediaFile pkg.TVShowFile, media pkg.TVEpisode, destination string) {
			log.Printf("Removing %s", source)
			pkg.AppendJobLog(fmt.Sprintf("Removing %s", source))
			err = os.Remove(source)
			if err != nil {
				log.Printf("Failed to remove %s : %s", source, err.Error())
				pkg.AppendJobLog(fmt.Sprintf("Failed to remove %s : %s", source, err.Error()))
			}
			// Upload destination to S3
			log.Printf("Uploading episode %d to S3...", media.ID)
			pkg.AppendJobLog(fmt.Sprintf("Uploading episode %d to S3...", media.ID))
			err = s.objectStorage.UploadMediaFiles(
				path.Join("tv-shows", strconv.Itoa(media.ID)),
				path.Join(destination, strconv.Itoa(media.ID)),
			)
			if err != nil {
				log.Printf("Failed to upload %s to S3 : %s", destination, err.Error())
				pkg.AppendJobLog(fmt.Sprintf("Failed to upload %s to S3 : %s", destination, err.Error()))
			}
			// Remove destination from local
			log.Printf("Removing %s from local storage", path.Join(destination, strconv.Itoa(media.ID)))
			pkg.AppendJobLog(fmt.Sprintf("Removing %s from local storage", path.Join(destination, strconv.Itoa(media.ID))))
			err = os.RemoveAll(path.Join(destination, strconv.Itoa(media.ID)))
			if err != nil {
				log.Printf("Failed to remove %s : %s", path.Join(destination, strconv.Itoa(media.ID)), err.Error())
				pkg.AppendJobLog(fmt.Sprintf("Failed to remove %s : %s", path.Join(destination, strconv.Itoa(media.ID)), err.Error()))
			}
		}(mediaFile, media, destination)
	}
	return nil
}
