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

func (m *MovieFile) String() string {
	return fmt.Sprintf("%-100s --> %s", m.Filename, m.SanitizedName)
}

type TVShowFile struct {
	Path          string
	SanitizedName string
	Season        int
	Episode       int
	Filename      string
	Extension     string
}

func (t TVShowFile) String() string {
	return fmt.Sprintf("%-100s --> %s S%.2dE%.2d", t.Filename, t.SanitizedName, t.Season, t.Episode)
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
	var mediaFiles = make([]MovieFile, 0)
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

// BuildTVShowTree recursively builds a tree of TV show files in the given source directory.
func BuildTVShowTree(source string) ([]TVShowFile, error) {
	// Read the directory entries from the source directory.
	entries, err := os.ReadDir(source)
	if err != nil {
		return nil, err
	}

	// Initialize a slice to store the TV show files.
	var tvShowFiles = make([]TVShowFile, 0)

	// Iterate over the entries in the source directory.
	for _, entry := range entries {
		// If the entry is a directory, recurse and add the resulting TV show files to the slice.
		if entry.IsDir() {
			recursiveTVShowFiles, err := BuildTVShowTree(filepath.Join(source, entry.Name()))
			if err != nil {
				return nil, err
			}
			tvShowFiles = append(tvShowFiles, recursiveTVShowFiles...)
		} else {
			// If the entry has an allowed extension, extract the TV show file information and add it to the slice.
			if hasAllowedExtension(entry.Name()) {
				var title, season, episode = SanitizeTVShowFilename(entry.Name())
				tvShowFile := TVShowFile{
					Path:          source,
					SanitizedName: title,
					Season:        season,
					Episode:       episode,
					Filename:      entry.Name(),
					Extension:     getExtension(entry.Name()),
				}
				tvShowFiles = append(tvShowFiles, tvShowFile)
			} else {
				log.Println("Not allowed extension: ", entry.Name())
			}
		}
	}

	return tvShowFiles, nil
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
