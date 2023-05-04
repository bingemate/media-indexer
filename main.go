package main

import (
	"flag"
	"github.com/bingemate/media-indexer/cmd"
	"github.com/bingemate/media-indexer/initializers"
	"log"
)

// @title Media Indexer API
// @description This is the API for the Media Indexer application
// @version 1.0
// @host localhost:8080
// @basePath /
func main() {
	var server = flag.Bool("serve", true, "Run server")
	flag.Parse()
	env, err := initializers.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}
	logFile := initializers.InitLog(env.LogFile)
	defer logFile.Close()
	if *server {
		log.Println("Starting server mode...")
		cmd.Serve(env)
	} else {
		log.Println("Starting cli mode...")
		cmd.ExecuteCli()
	}
}
