package cmd

import (
	"github.com/bingemate/media-indexer/internal/features"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "media-indexer",
	Short: "Media Indexer",
	Long:  "Media Indexer",
	Run: func(cmd *cobra.Command, args []string) {
		source, _ := cmd.Flags().GetString("source")
		destination, _ := cmd.Flags().GetString("destination")
		tmdbApiKey, _ := cmd.Flags().GetString("tmdb-api-key")
		if strings.TrimSpace(source) == "" || strings.TrimSpace(destination) == "" {
			log.Println("Source and destination are required")
			_ = cmd.Help()
			os.Exit(1)
		}
		main(source, destination, tmdbApiKey)
	},
}

func ExecuteCli() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringP("source", "s", "", "Source directory")
	rootCmd.Flags().StringP("destination", "d", "", "Destination directory")
	rootCmd.Flags().StringP("tmdb-api-key", "t", "", "TMDB API Key")
}

func main(source, destination, tmdbApiKey string) {
	log.Printf("Source: %s\n", source)
	log.Printf("Destination: %s\n", destination)
	var movieScanner = features.NewMovieScanner(source, destination, tmdbApiKey)
	_, err := movieScanner.ScanMovies()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done")
}
