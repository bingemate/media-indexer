package pkg

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

var deleteRegexes = []*regexp.Regexp{
	regexp.MustCompile(`[\[\(].*?[\]\)]|-\s*\d+p.*`), // regex pour supprimer les informations non pertinentes
	regexp.MustCompile(`\s+$`),                       // regex pour supprimer les espaces en fin de chaîne
}

var spaceRegexes = []*regexp.Regexp{
	regexp.MustCompile(`[\W_]+`), // regex pour supprimer les caractères spéciaux
}

var nonASCIIRegex = regexp.MustCompile(`[^\x00-\x7F]+`) // regex pour remplacer les caractères non ASCII

var extractDateRegex = regexp.MustCompile(`^(.+?)(\d{4}?).*(\d{4}.*)?$`) // Expression régulière pour extraire le nom et l'année du fichier

var tvShowRegex = regexp.MustCompile(`^(.+?)(?:[sS])?(\d{1,})?(?:[eExX])?(\d{2,})(?:.*|$)`) // regex to extract title, season number, and episode number

// SanitizeMovieFilename sanitize a movie filename by removing non ASCII characters, removing non-relevant information, returning the name and the year
func SanitizeMovieFilename(filename string) (string, string) {
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	for _, regex := range spaceRegexes {
		filename = regex.ReplaceAllString(filename, " ")
	}
	filename = nonASCIIRegex.ReplaceAllStringFunc(filename, func(s string) string {
		return strings.ToLower(strings.TrimFunc(s, func(r rune) bool {
			return !unicode.IsLetter(r)
		}))
	})

	for _, regex := range deleteRegexes {
		filename = regex.ReplaceAllString(filename, "")
	}

	matches := extractDateRegex.FindStringSubmatch(filename)
	if len(matches) < 3 {
		return filename, ""
	}
	name := strings.TrimSpace(matches[1])
	year := strings.TrimSpace(matches[2])
	return name, year
}

// SanitizeTVShowFilename separates a TV show filename into title, season number, and episode number
func SanitizeTVShowFilename(filename string) (string, string, string) {
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	for _, regex := range spaceRegexes {
		filename = regex.ReplaceAllString(filename, " ")
	}
	filename = nonASCIIRegex.ReplaceAllStringFunc(filename, func(s string) string {
		return strings.ToLower(strings.TrimFunc(s, func(r rune) bool {
			return !unicode.IsLetter(r)
		}))
	})

	for _, regex := range deleteRegexes {
		filename = regex.ReplaceAllString(filename, "")
	}

	matches := tvShowRegex.FindStringSubmatch(filename)
	if len(matches) < 3 {
		return filename, "", ""
	}

	title := strings.TrimSpace(matches[1])
	seasonNumberStr := matches[2]
	episodeNumberStr := matches[3]

	// Default to season 1 if no season number is specified
	seasonNumber := "1"
	if seasonNumberStr != "" {
		seasonNumber = seasonNumberStr
	}

	return title, seasonNumber, episodeNumberStr
}
