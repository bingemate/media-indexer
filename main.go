package main

import (
	"flag"
	"fmt"
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
		cmd.Serve(env)
	} else {
		fmt.Println("Running CLI")
		//cmd.ExecuteCli()
	}
}
