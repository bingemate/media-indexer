package main

import (
	"flag"
	"github.com/bingemate/media-indexer/cmd"
	"github.com/bingemate/media-indexer/initializers"
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
	if *server {
		log.Println("Starting server mode...")
		cmd.Serve(env)
	} else {
		log.Println("Starting cli mode...")
		cmd.ExecuteCli()
	}
}
