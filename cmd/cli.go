package cmd

import (
	"github.com/bingemate/media-indexer/internal"
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
		if strings.TrimSpace(source) == "" || strings.TrimSpace(destination) == "" {
			log.Println("Source and destination are required")
			_ = cmd.Help()
			os.Exit(1)
		}
		main(source, destination)
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
}

func main(source, destination string) {
	log.Printf("Source: %s\n", source)
	log.Printf("Destination: %s\n", destination)
	err := internal.ScanMovieFolder(source, destination)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done")
}
