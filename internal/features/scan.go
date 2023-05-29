package features

import (
	"errors"
	"github.com/bingemate/media-indexer/internal/repository"
	"github.com/bingemate/media-indexer/pkg"
	"log"
	"path"
	"sync"
)

var (
	scannerMutex = &sync.Mutex{}
)

// MovieScanner represents a struct that scans movie folders to search for movie files and move them.
type MovieScanner struct {
	source          string                      // Source directory path to scan for movies.
	destination     string                      // Destination directory path to move the found movies.
	mediaClient     pkg.MediaClient             // Media client object to search for movies on TMDB.
	mediaRepository *repository.MediaRepository // Media repository object to save the media files and their details.
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
}

// TVScannerResult represents a struct that holds the results of scanning and moving TV show files.
type TVScannerResult struct {
	Source    string        // Source filename.
	TVEpisode pkg.TVEpisode // TV episode details returned by TMDB.
}

// NewMovieScanner returns a new instance of MovieScanner with given source directory, target directory, and TMDB API key.
func NewMovieScanner(source, destination string, mediaClient pkg.MediaClient, mediaRepository *repository.MediaRepository) *MovieScanner {
	return &MovieScanner{
		source:          source,
		destination:     destination,
		mediaClient:     mediaClient,
		mediaRepository: mediaRepository,
	}
}

// NewTVScanner returns a new instance of TVScanner with given source directory, target directory, and TMDB API key.
func NewTVScanner(source, destination string, mediaClient pkg.MediaClient, mediaRepository *repository.MediaRepository) *TVScanner {
	return &TVScanner{
		source:          source,
		destination:     destination,
		mediaClient:     mediaClient,
		mediaRepository: mediaRepository,
	}
}

// ScanMovies scans the source directory for movies and moves them to the destination directory.
// It returns a slice of MovieScannerResult and an error if there is any.
func (s *MovieScanner) ScanMovies() (*[]MovieScannerResult, error) {
	// Locks the scanner to prevent concurrent scanning
	locked := scannerMutex.TryLock()
	if !locked {
		return nil, errors.New("scanner is currently running")
	}
	defer scannerMutex.Unlock()

	uploadLocked := uploadLock.TryLock()
	if !uploadLocked {
		return nil, errors.New("upload is currently running")
	}
	defer uploadLock.Unlock()
	mediaFiles, err := s.scanMovieFolder()
	if err != nil {
		return nil, err
	}
	atomicMovieList := s.retrieveMovieList(mediaFiles)

	result := s.buildMovieScannerResult(atomicMovieList)

	//log.Printf("Moving %d movies to %s...", len(*result), s.destination)

	// Moves the movies to the destination directory and returns an error if it fails
	err = s.processMovies(atomicMovieList, s.destination)
	if err != nil {
		log.Printf("Failed to move movies to %s: %v", s.destination, err)
		return nil, err
	}

	log.Printf("Successfully moved %d movies to %s.", len(*result), s.destination)

	err = pkg.ClearFolderContent(s.source)
	if err != nil {
		log.Printf("Failed to clear source folder content: %v", err)
		return nil, err
	}

	// Returns the result slice and no error
	return result, nil
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

	// Builds the directory tree from the source directory and returns an error if it fails
	mediaFiles, err := pkg.BuildMovieTree(s.source)
	if err != nil {
		log.Printf("Failed to scan source tree: %v", err)
		return nil, err
	}

	log.Printf("Scanning %d files in %s...", len(mediaFiles), s.source)
	return &mediaFiles, nil
}

func (s *MovieScanner) retrieveMovieList(mediaFiles *[]pkg.MovieFile) *pkg.AtomicMovieList {
	// Initializes a WaitGroup and an AtomicMovieList
	var wg sync.WaitGroup
	var atomicMovieList = pkg.NewAtomicMovieList()
	wg.Add(len(*mediaFiles))

	// Iterates through each media file and spawns a goroutine to search its information
	for _, mediaFile := range *mediaFiles {
		go func(mediaFile pkg.MovieFile) {
			defer wg.Done()

			log.Printf("Searching for movie information for file %s...", mediaFile.Filename)

			// Searches for the movie information for the current media file and adds the result to the AtomicMovieList
			media, ok := searchMovie(&mediaFile, s.mediaClient)
			if !ok {
				log.Printf("Failed to find movie information for file %s.", mediaFile.Filename)
				return
			}
			atomicMovieList.LinkMediaFile(mediaFile, media)
		}(mediaFile)
	}

	// Waits for all goroutines to finish
	wg.Wait()

	log.Println("Movie scan complete.")
	return atomicMovieList
}

// ScanTV scans the source directory for TV shows and moves them to the destination directory.
// It returns a slice of TVScannerResult and an error if there is any.
func (s *TVScanner) ScanTV() (*[]TVScannerResult, error) {
	// Locks the scanner to prevent concurrent scanning
	locked := scannerMutex.TryLock()
	if !locked {
		return nil, errors.New("scanner is currently running")
	}
	defer scannerMutex.Unlock()

	uploadLocked := uploadLock.TryLock()
	if !uploadLocked {
		return nil, errors.New("upload is currently running")
	}
	defer uploadLock.Unlock()
	mediaFiles, err := s.scanTVFolder()
	if err != nil {
		return nil, err
	}

	atomicMediaList := s.retrieveTvList(mediaFiles)

	result := s.buildTVScannerResult(atomicMediaList)

	//log.Printf("Moving %d TV shows to %s...", len(*result), s.destination)

	// Moves the TV shows to the destination directory and returns an error if it fails
	err = s.processTVEpisodes(atomicMediaList, s.destination)
	if err != nil {
		log.Printf("Failed to move TV shows to %s: %v", s.destination, err)
		return nil, err
	}

	log.Printf("Successfully moved %d TV shows to %s.", len(*result), s.destination)

	err = pkg.ClearFolderContent(s.source)
	if err != nil {
		log.Printf("Failed to clear source folder content: %v", err)
		return nil, err
	}

	// Returns the result slice and no error
	return result, nil
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
	// Initializes a WaitGroup and an AtomicMovieList
	var wg sync.WaitGroup
	var atomicMediaList = pkg.NewAtomicTVEpisodeList()
	wg.Add(len(*mediaFiles))

	// Iterates through each media file and spawns a goroutine to search its information
	for _, mediaFile := range *mediaFiles {
		go func(mediaFile pkg.TVShowFile) {
			defer wg.Done()

			// Logs that the function is searching for TV show information for the current file
			log.Printf("Searching for TV show information for file %s...", mediaFile.Filename)

			// Searches for the TV show information for the current media file and adds the result to the AtomicMovieList
			media, ok := searchTVEpisode(&mediaFile, s.mediaClient)
			if !ok {
				log.Printf("Failed to find TV show information for file %s.", mediaFile.Filename)
				return
			}
			log.Printf("Found TV show information for file %s:", mediaFile.Filename)
			log.Println(media)
			atomicMediaList.LinkMediaFile(mediaFile, media)
		}(mediaFile)
	}

	// Waits for all goroutines to finish
	wg.Wait()

	// Logs that the TV show scan is complete
	log.Println("TV show scan complete.")
	return atomicMediaList
}

func (s *TVScanner) scanTVFolder() (*[]pkg.TVShowFile, error) {
	// Logs that the function is scanning the source directory for TV shows
	log.Printf("Scanning %s for TV shows...", s.source)

	// Builds the directory tree from the source directory and returns an error if it fails
	mediaFiles, err := pkg.BuildTVShowTree(s.source)
	if err != nil {
		log.Printf("Failed to scan source tree: %v", err)
		return nil, err
	}

	log.Printf("Scanning %d files in %s...", len(mediaFiles), s.source)
	return &mediaFiles, nil
}

// searchMovie searches for a movie on TMDB using the media file name and year, returning the movie details and a boolean indicating whether it was found.
func searchMovie(mediaFile *pkg.MovieFile, client pkg.MediaClient) (pkg.Movie, bool) {
	result, err := client.SearchMovie(mediaFile.SanitizedName, mediaFile.Year)
	if err != nil {
		log.Printf("Error while media search on %s : %s. Sanitized name was : %s", mediaFile.Filename, err.Error(), mediaFile.SanitizedName)
		return pkg.Movie{}, false
	}
	return result, true
}

// searchTVEpisode searches for a TV show on TMDB using the media file name and year, returning the TV show details and a boolean indicating whether it was found.
func searchTVEpisode(mediaFile *pkg.TVShowFile, client pkg.MediaClient) (pkg.TVEpisode, bool) {
	result, err := client.SearchTVShow(mediaFile.SanitizedName, mediaFile.Season, mediaFile.Episode)
	if err != nil {
		log.Printf("Error while media search on %s : %s. Sanitized name was : %s", mediaFile.Filename, err.Error(), mediaFile.SanitizedName)
		return pkg.TVEpisode{}, false
	}
	return result, true
}

// processMovies moves the media files to the destination directory path provided as argument.
// It returns an error if the destination directory does not exist or if there was an error while moving the file.
func (s *MovieScanner) processMovies(movieList *pkg.AtomicMovieList, destination string) error {
	if !pkg.IsDirectoryExists(destination) {
		return errors.New("destination directory does not exists")
	}
	for mediaFile, media := range movieList.GetAll() {
		var source = path.Join(mediaFile.Path, mediaFile.Filename)
		err := s.mediaRepository.IndexMovie(media, source, s.destination)
		if err != nil {
			return err
		}
		log.Printf("Processed %-60s - %s %s", mediaFile.Filename, media.Name, media.Year())
	}
	return nil
}

// processTVEpisodes moves the media files to the destination directory path provided as argument.
// It returns an error if the destination directory does not exist or if there was an error while moving the file.
func (s *TVScanner) processTVEpisodes(tvList *pkg.AtomicTVEpisodeList, destination string) error {
	if !pkg.IsDirectoryExists(destination) {
		return errors.New("destination directory does not exists")
	}

	for mediaFile, media := range tvList.GetAll() {
		var source = path.Join(mediaFile.Path, mediaFile.Filename)
		err := s.mediaRepository.IndexTvEpisode(media, source, s.destination)
		if err != nil {
			return err
		}
		log.Printf("Processed %-60s - %s %s", mediaFile.Filename, media.Name, media.Year())
	}
	return nil
}
