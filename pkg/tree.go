package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type MovieFile struct {
	Path          string
	Filename      string
	SanitizedName string
	Year          string
	Extension     string
}

func (m MovieFile) String() string {
	return fmt.Sprintf("%-100s --> %s", m.Filename, m.SanitizedName)
}

var allowedExtension = []string{
	".mp4",
	".mkv",
	".avi",
}

func BuildMovieTree(source string) ([]MovieFile, error) {
	entries, err := os.ReadDir(source)
	if err != nil {
		return nil, err
	}
	var mediaFiles []MovieFile
	for _, entry := range entries {
		if entry.IsDir() {
			recursiveMediaFiles, err := BuildMovieTree(filepath.Join(source, entry.Name()))
			if err != nil {
				return nil, err
			}
			mediaFiles = append(mediaFiles, recursiveMediaFiles...)
		} else {
			if hasAllowedExtension(entry.Name()) {
				var title, year = SanitizeMovieFilename(entry.Name())
				mediaFile := MovieFile{
					Path:          source,
					Filename:      entry.Name(),
					SanitizedName: title,
					Year:          year,
					Extension:     getExtension(entry.Name()),
				}
				mediaFiles = append(mediaFiles, mediaFile)
			} else {
				log.Println("Not allowed extension: ", entry.Name())
			}
		}
	}
	return mediaFiles, nil
}

func getExtension(name string) string {
	return name[strings.LastIndex(name, "."):]
}

func hasAllowedExtension(filename string) bool {
	for _, extension := range allowedExtension {
		if strings.HasSuffix(filename, extension) {
			return true
		}
	}
	return false
}
