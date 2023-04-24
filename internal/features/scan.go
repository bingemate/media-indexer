package features

import (
	"errors"
	"fmt"
	"github.com/bingemate/media-indexer/internal/repository"
	"github.com/bingemate/media-indexer/pkg"
	"log"
	"path"
	"path/filepath"
	"sync"
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
	Source      string    // Source filename.
	Destination string    // Full destination path of the moved file.
	Movie       pkg.Movie // Movie details returned by TMDB.
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
	Source      string        // Source filename.
	Destination string        // Full destination path of the moved file.
	TVEpisode   pkg.TVEpisode // TV episode details returned by TMDB.
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
	mediaFiles, err := s.scanMovieFolder()
	if err != nil {
		return nil, err
	}

	atomicMovieList := s.retrieveMovieList(mediaFiles)

	result := s.buildMovieScannerResult(atomicMovieList)

	log.Printf("Moving %d movies to %s...", len(*result), s.destination)

	// Moves the movies to the destination directory and returns an error if it fails
	err = s.moveMovies(atomicMovieList, s.destination)
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
			Source:      mediaFile.Filename,
			Destination: buildMovieFilename(media, mediaFile.Extension),
			Movie:       media,
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
	mediaFiles, err := s.scanTVFolder()
	if err != nil {
		return nil, err
	}

	atomicMediaList := s.retrieveTvList(mediaFiles)

	result := s.buildTVScannerResult(atomicMediaList)

	log.Printf("Moving %d TV shows to %s...", len(*result), s.destination)

	// Moves the TV shows to the destination directory and returns an error if it fails
	err = s.moveTVEpisodes(atomicMediaList, s.destination)
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
			Source:      mediaFile.Filename,
			Destination: filepath.Join(media.Name, buildTVEpisodeFilename(media, mediaFile.Extension)),
			TVEpisode:   media,
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

// moveMovies moves the media files to the destination directory path provided as argument.
// It returns an error if the destination directory does not exist or if there was an error while moving the file.
func (s *MovieScanner) moveMovies(movieList *pkg.AtomicMovieList, destination string) error {
	if !pkg.IsDirectoryExists(destination) {
		return errors.New("destination directory does not exists")
	}
	for mediaFile, media := range movieList.GetAll() {
		var source = path.Join(mediaFile.Path, mediaFile.Filename)
		var movieFilename = buildMovieFilename(media, mediaFile.Extension)
		var destination = path.Join(
			destination,
			movieFilename,
		)
		err := pkg.MoveFile(source, destination)
		if err != nil {
			return err
		}
		err = s.mediaRepository.IndexMovie(&media, s.destination, movieFilename)
		if err != nil {
			return err
		}
		log.Printf("Processed %-60s - %s %s", mediaFile.Filename, media.Name, media.Year())
	}
	return nil
}

// moveTVEpisodes moves the media files to the destination directory path provided as argument.
// It returns an error if the destination directory does not exist or if there was an error while moving the file.
func (s *TVScanner) moveTVEpisodes(tvList *pkg.AtomicTVEpisodeList, destination string) error {
	if !pkg.IsDirectoryExists(destination) {
		return errors.New("destination directory does not exists")
	}
	for mediaFile, media := range tvList.GetAll() {
		var source = path.Join(mediaFile.Path, mediaFile.Filename)
		var destination = path.Join(
			destination,
			media.Name,
			buildTVEpisodeFilename(media, mediaFile.Extension),
		)
		err := pkg.MoveFile(source, destination)
		if err != nil {
			return err
		}
		log.Printf("Processed %-60s - %s %s", mediaFile.Filename, media.Name, media.Year())
	}
	return nil
}

// buildTVEpisodeFilename builds a TV show filename using the TV show name, season and episode number.
func buildTVEpisodeFilename(media pkg.TVEpisode, extension string) string {
	return fmt.Sprintf("%s - S%dE%.2d%s", media.Name, media.Season, media.Episode, extension)
}

// buildMovieFilename builds a movie filename using the movie name and year.
func buildMovieFilename(media pkg.Movie, extension string) string {
	return fmt.Sprintf("%s - %s%s", media.Name, media.Year(), extension)
}
