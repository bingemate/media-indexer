package main

import (
	"flag"
	"github.com/bingemate/media-indexer/cmd"
	"github.com/bingemate/media-indexer/initializers"
	"github.com/bingemate/media-indexer/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

func main() {
	var server = flag.Bool("serve", false, "Run server")
	flag.Parse()
	env, err := initializers.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}
	logFile := initializers.InitLog(env.LogFile)
	defer logFile.Close()

	//test(db)
	if *server {
		log.Println("Starting server mode...")
		cmd.Serve(env)
	} else {
		log.Println("Starting cli mode...")
		cmd.ExecuteCli()
	}
}

func test(db *gorm.DB) {
	var mediaFile = repository.MediaFile{
		Filename: "test.mkv",
		Codec:    repository.VideoCodecH264,
		Size:     1266565656.58,
		Duration: 4600.58,
		Subtitles: []repository.Subtitle{
			{
				Codec:    repository.SubtitleCodecSRT,
				Language: "français",
			},
		},
	}

	/*	var subtitle = repository.Subtitles{
		Codec:    repository.SRT,
		Language: "français",
	}*/

	db.Create(&mediaFile)
	var mediaFiles []repository.MediaFile
	db.Joins(clause.Associations).Find(&mediaFiles)
	log.Println(mediaFiles)
}
