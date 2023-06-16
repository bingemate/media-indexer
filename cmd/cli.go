package cmd

import (
	"github.com/bingemate/media-indexer/initializers"
	"github.com/bingemate/media-indexer/internal/features"
	"github.com/bingemate/media-indexer/internal/repository"
	"github.com/bingemate/media-indexer/pkg"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "media-indexer",
	Short: "Media Indexer",
	Long:  "Media Indexer",
	Run: func(cmd *cobra.Command, args []string) {
		//source, _ := cmd.Flags().GetString("source")
		//destination, _ := cmd.Flags().GetString("destination")
		//tmdbApiKey, _ := cmd.Flags().GetString("tmdb-api-key")
		//if strings.TrimSpace(source) == "" || strings.TrimSpace(destination) == "" {
		//	log.Println("Source and destination are required")
		//	_ = cmd.Help()
		//	os.Exit(1)
		//}
		env, err := initializers.LoadEnv()
		if err != nil {
			log.Fatal(err)
		}

		main(env)
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

func main(env initializers.Env) {
	log.Printf("Source: %s\n", env.MovieSourceFolder)
	log.Printf("Destination: %s\n", env.MovieTargetFolder)
	var mediaClient = pkg.NewMediaClient(env.TMDBApiKey)
	db, err := initializers.ConnectToDB(env)
	if err != nil {
		log.Fatal(err)
	}
	var mediaRepository = repository.NewMediaRepository(db)
	var movieScanner = features.NewMovieScanner(env.MovieSourceFolder, env.MovieTargetFolder, mediaClient, mediaRepository)
	err = movieScanner.ScanMovies()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done")
}
