package internal

import (
	"errors"
	"fmt"
	"github.com/bingemate/media-indexer/pkg"
	"log"
	"sync"
)

func ScanMovieFolder(source, destination, tmdbApiKey string) error {
	sourceTree, err := pkg.BuildTree(source)
	if err != nil {
		return err
	}
	var client = pkg.NewMediaClient(tmdbApiKey)
	var wg sync.WaitGroup
	var atomicMediaList = pkg.NewAtomicMediaList()
	wg.Add(len(sourceTree))

	for _, mediaFile := range sourceTree {
		go func(mediaFile pkg.MediaFile) {
			defer wg.Done()
			Media, ok := searchMedia(&mediaFile, client)
			if !ok {
				return
			}
			atomicMediaList.LinkMediaFile(mediaFile, Media)
		}(mediaFile)
	}

	wg.Wait()
	return moveMedias(atomicMediaList.GetAll(), destination)
}

func searchMedia(mediaFile *pkg.MediaFile, client pkg.MediaClient) (pkg.Media, bool) {
	result, err := client.SearchMovie(mediaFile.SanitizedName, mediaFile.Year)
	if err != nil {
		log.Printf("Error while media search on %s : %s. Sanitized name was : %s", mediaFile.Filename, err.Error(), mediaFile.SanitizedName)
		return pkg.Media{}, false
	}
	return result, true
}

func moveMedias(mediaList map[pkg.MediaFile]pkg.Media, destination string) error {
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
		// TODO: Index in database
		log.Printf("Processed %-60s - %s %s", mediaFile.Filename, media.Name, media.Year())
	}
	return nil
}
