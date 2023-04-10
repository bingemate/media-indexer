package tree

import (
	"fmt"
	"github.com/bingemate/media-indexer/pkg/sanitizer"
	"os"
	"strings"
)

type MediaFile struct {
	Path          string
	Filename      string
	SanitizedName string
}

func (m MediaFile) String() string {
	return fmt.Sprintf("%-100s --> %s", m.Filename, m.SanitizedName)
}

var allowedExtension = []string{
	".mp4",
	".mkv",
}

func BuildTree(source string) ([]MediaFile, error) {
	entries, err := os.ReadDir(source)
	if err != nil {
		return nil, err
	}
	var mediaFiles []MediaFile
	for _, entry := range entries {
		if entry.IsDir() {
			recursiveMediaFiles, err := BuildTree(source + "/" + entry.Name())
			if err != nil {
				return nil, err
			}
			mediaFiles = append(mediaFiles, recursiveMediaFiles...)
		} else {
			for _, extension := range allowedExtension {
				if strings.HasSuffix(entry.Name(), extension) {
					mediaFile := MediaFile{
						Path:          source,
						Filename:      entry.Name(),
						SanitizedName: sanitizer.SanitizeFilename(entry.Name()),
					}
					mediaFiles = append(mediaFiles, mediaFile)
				}
			}
		}
	}
	return mediaFiles, nil
}
