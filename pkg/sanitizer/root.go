package sanitizer

import (
	"fmt"
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

// var extractDateRegex = regexp.MustCompile(`^(.+?)(?:\s*\((\d{4})\))?[\s\._-]*[^\\\/]*$`) // Expression régulière pour extraire le nom et l'année du fichier
var extractDateRegex = regexp.MustCompile(`^(.+?)(\d{4}?).*(\d{4}.*)?$`) // Expression régulière pour extraire le nom et l'année du fichier

func SanitizeFilename(filename string) string {
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
		return filename
	}
	name := strings.TrimSpace(matches[1])
	year := strings.TrimSpace(matches[2])

	return fmt.Sprintf("%s %s", name, year)
}
