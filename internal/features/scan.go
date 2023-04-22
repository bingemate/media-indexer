package features

import (
	"errors"
	"fmt"
	"github.com/bingemate/media-indexer/pkg"
	"log"
	"sync"
)

// MovieScanner represents a struct that scans movie folders to search for movie files and move them.
type MovieScanner struct {
	source      string          // Source directory path to scan for movies.
	destination string          // Destination directory path to move the found movies.
	mediaClient pkg.MediaClient // Media client object to search for movies on TMDB.
}

// MovieScannerResult represents a struct that holds the results of scanning and moving movie files.
type MovieScannerResult struct {
	Source      string    // Source filename.
	Destination string    // Full destination path of the moved file.
	Movie       pkg.Movie // Movie details returned by TMDB.
}

// NewMovieScanner returns a new instance of MovieScanner with given source directory, target directory, and TMDB API key.
func NewMovieScanner(source, destination, tmdbAPIKey string) *MovieScanner {
	return &MovieScanner{
		source:      source,
		destination: destination,
		mediaClient: pkg.NewMediaClient(tmdbAPIKey),
	}
}

// ScanMovieFolder scans the source directory for movies and moves them to the destination directory.
// It returns a slice of MovieScannerResult and an error if there is any.
func (s *MovieScanner) ScanMovieFolder() ([]MovieScannerResult, error) {
	// Logs that the function is scanning the source directory for movies
	log.Printf("Scanning %s for movies...", s.source)

	// Builds the directory tree from the source directory and returns an error if it fails
	sourceTree, err := pkg.BuildTree(s.source)
	if err != nil {
		log.Printf("Failed to scan source tree: %v", err)
		return nil, err
	}

	// Logs the number of files found in the source directory
	log.Printf("Scanning %d files in %s...", len(sourceTree), s.source)

	// Initializes a WaitGroup and an AtomicMediaList
	var wg sync.WaitGroup
	var atomicMediaList = pkg.NewAtomicMediaList()
	wg.Add(len(sourceTree))

	// Iterates through each media file and spawns a goroutine to search its information
	for _, mediaFile := range sourceTree {
		go func(mediaFile pkg.MediaFile) {
			defer wg.Done()

			// Logs that the function is searching for movie information for the current file
			log.Printf("Searching for movie information for file %s...", mediaFile.Filename)

			// Searches for the movie information for the current media file and adds the result to the AtomicMediaList
			Media, ok := searchMovie(&mediaFile, s.mediaClient)
			if !ok {
				return
			}
			atomicMediaList.LinkMediaFile(mediaFile, Media)
		}(mediaFile)
	}

	// Waits for all goroutines to finish
	wg.Wait()

	// Logs that the movie scan is complete
	log.Println("Movie scan complete.")

	// Initializes an empty slice of MovieScannerResult
	var result = make([]MovieScannerResult, 0)

	// Iterates through each media file and its corresponding movie information in the AtomicMediaList and adds it to the result slice
	for mediaFile, media := range atomicMediaList.GetAll() {
		result = append(result, MovieScannerResult{
			Source:      mediaFile.Filename,
			Destination: fmt.Sprintf("%s - %s%s", media.Name, media.Year(), mediaFile.Extension),
			Movie:       media,
		})
	}

	// Logs that it's moving the movies to the destination directory
	log.Printf("Moving %d movies to %s...", len(result), s.destination)

	// Moves the movies to the destination directory and returns an error if it fails
	err = moveMovies(atomicMediaList.GetAll(), s.destination)
	if err != nil {
		log.Printf("Failed to move movies to %s: %v", s.destination, err)
		return nil, err
	}

	// Logs that it has successfully moved the movies to the destination directory
	log.Printf("Successfully moved %d movies to %s.", len(result), s.destination)

	// Returns the result slice and no error
	return result, nil
}

// searchMovie searches for a movie on TMDB using the media file name and year, returning the movie details and a boolean indicating whether it was found.
func searchMovie(mediaFile *pkg.MediaFile, client pkg.MediaClient) (pkg.Movie, bool) {
	result, err := client.SearchMovie(mediaFile.SanitizedName, mediaFile.Year)
	if err != nil {
		log.Printf("Error while media search on %s : %s. Sanitized name was : %s", mediaFile.Filename, err.Error(), mediaFile.SanitizedName)
		return pkg.Movie{}, false
	}
	return result, true
}

// moveMovies moves the media files to the destination directory path provided as argument.
// It returns an error if the destination directory does not exist or if there was an error while moving the file.
func moveMovies(mediaList map[pkg.MediaFile]pkg.Movie, destination string) error {
	if !pkg.IsDirectoryExists(destination) {
		return errors.New("destination directory does not exists")
	}
	for mediaFile, media := range mediaList {
		var source = fmt.Sprintf("%s/%s", mediaFile.Path, mediaFile.Filename)
		var destination = fmt.Sprintf("%s/%s - %s%s", destination, media.Name, media.Year(), mediaFile.Extension)
		err := pkg.MoveFile(source, destination)
		if err != nil {
			return err
		}
		log.Printf("Processed %-60s - %s %s", mediaFile.Filename, media.Name, media.Year())
	}
	return nil
}
