package pkg

import (
	"fmt"
	"os"
	"strings"
)

type MediaFile struct {
	Path          string
	Filename      string
	SanitizedName string
	Year          string
	Extension     string
}

func (m MediaFile) String() string {
	return fmt.Sprintf("%-100s --> %s", m.Filename, m.SanitizedName)
}

var allowedExtension = []string{
	".mp4",
	".mkv",
	".avi",
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
			if hasAllowedExtension(entry.Name()) {
				var title, year = SanitizeFilename(entry.Name())
				mediaFile := MediaFile{
					Path:          source,
					Filename:      entry.Name(),
					SanitizedName: title,
					Year:          year,
					Extension:     getExtension(entry.Name()),
				}
				mediaFiles = append(mediaFiles, mediaFile)
			} else {
				fmt.Println("Not allowed extension: ", entry.Name())
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
