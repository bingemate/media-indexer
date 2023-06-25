package pkg

import (
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var deleteRegexes = []*regexp.Regexp{
	regexp.MustCompile(`[\[\(].*?[\]\)]|-\s*\d+p.*`), // regex pour supprimer les informations non pertinentes
	regexp.MustCompile(`\s+$`),                       // regex pour supprimer les espaces en fin de chaîne
}

var spaceRegexes = []*regexp.Regexp{
	regexp.MustCompile(`[^\pL\s_]+`), // regex pour supprimer les caractères spéciaux
}

var extractDateRegex = regexp.MustCompile(`^(.+?)(\d{4}?).*(\d{4}.*)?$`) // Expression régulière pour extraire le nom et l'année du fichier

var tvShowRegex = regexp.MustCompile(`^(.+?)(?:[sS])?(\d{1,})?(?:[eExX])?(\d{2,})(?:.*|$)`) // regex to extract title, season number, and episode number

var isMn = func(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// SanitizeMovieFilename sanitize a movie filename by removing non ASCII characters, removing non-relevant information, returning the name and the year
func SanitizeMovieFilename(filename string) (string, string) {
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	for _, regex := range spaceRegexes {
		filename = regex.ReplaceAllString(filename, " ")
	}
	filename = removeAccents(filename)

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
func SanitizeTVShowFilename(filename string) (string, int, int) {
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	for _, regex := range spaceRegexes {
		filename = regex.ReplaceAllString(filename, " ")
	}
	filename = removeAccents(filename)

	for _, regex := range deleteRegexes {
		filename = regex.ReplaceAllString(filename, "")
	}

	matches := tvShowRegex.FindStringSubmatch(filename)
	if len(matches) < 4 {
		return filename, 0, 0
	}

	title := strings.TrimSpace(matches[1])
	seasonNumberStr := matches[2]
	episodeNumberStr := matches[3]

	episodeNumber, err := strconv.Atoi(episodeNumberStr)
	if err != nil {
		log.Println("Error parsing episode number: ", err)
		episodeNumber = 0
	}

	// Default to season 1 if no season number is specified
	seasonNumber := 1
	if seasonNumberStr != "" {
		seasonNumber, err = strconv.Atoi(seasonNumberStr)
		if err != nil {
			log.Println("Error parsing season number: ", err)
			seasonNumber = 1
		}
	}

	return title, seasonNumber, episodeNumber
}
