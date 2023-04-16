package features

import (
	"errors"
	"fmt"
	"github.com/bingemate/media-indexer/pkg"
	"log"
	"sync"
)

type MovieScanner struct {
	source      string
	destination string
	mediaClient pkg.MediaClient
}

type MovieScannerResult struct {
	Source      string
	Destination string
	Movie       pkg.Movie
}

func NewMovieScanner(source, destination, tmdbApiKey string) *MovieScanner {
	return &MovieScanner{
		source:      source,
		destination: destination,
		mediaClient: pkg.NewMediaClient(tmdbApiKey),
	}
}

func (s *MovieScanner) ScanMovieFolder() ([]MovieScannerResult, error) {
	log.Printf("Scanning %s for movies...", s.source)
	sourceTree, err := pkg.BuildTree(s.source)
	if err != nil {
		log.Printf("Failed to scan source tree: %v", err)
		return nil, err
	}

	log.Printf("Scanning %d files in %s...", len(sourceTree), s.source)

	var wg sync.WaitGroup
	var atomicMediaList = pkg.NewAtomicMediaList()
	wg.Add(len(sourceTree))

	for _, mediaFile := range sourceTree {
		go func(mediaFile pkg.MediaFile) {
			defer wg.Done()

			log.Printf("Searching for movie information for file %s...", mediaFile.Filename)

			Media, ok := searchMovie(&mediaFile, s.mediaClient)
			if !ok {
				return
			}
			atomicMediaList.LinkMediaFile(mediaFile, Media)
		}(mediaFile)
	}

	wg.Wait()

	log.Println("Movie scan complete.")

	var result = make([]MovieScannerResult, 0)
	for mediaFile, media := range atomicMediaList.GetAll() {
		result = append(result, MovieScannerResult{
			Source:      mediaFile.Filename,
			Destination: fmt.Sprintf("%s - %s%s", media.Name, media.Year(), mediaFile.Extension),
			Movie:       media,
		})
	}

	log.Printf("Moving %d movies to %s...", len(result), s.destination)

	err = moveMovies(atomicMediaList.GetAll(), s.destination)
	if err != nil {
		log.Printf("Failed to move movies to %s: %v", s.destination, err)
		return nil, err
	}

	log.Printf("Successfully moved %d movies to %s.", len(result), s.destination)

	return result, nil
}

func searchMovie(mediaFile *pkg.MediaFile, client pkg.MediaClient) (pkg.Movie, bool) {
	result, err := client.SearchMovie(mediaFile.SanitizedName, mediaFile.Year)
	if err != nil {
		log.Printf("Error while media search on %s : %s. Sanitized name was : %s", mediaFile.Filename, err.Error(), mediaFile.SanitizedName)
		return pkg.Movie{}, false
	}
	return result, true
}

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
