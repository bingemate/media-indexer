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
	db, err := initializers.ConnectToDB(env)
	if err != nil {
		log.Fatal(err)
	}
	test(db)
	if *server {
		log.Println("Starting server mode...")
		cmd.Serve(env)
	} else {
		log.Println("Starting cli mode...")
		cmd.ExecuteCli()
	}
}
