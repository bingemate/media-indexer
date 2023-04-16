package main

import (
	"flag"
	"fmt"
	"github.com/bingemate/media-indexer/cmd"
	"github.com/bingemate/media-indexer/initializers"
	"io"
	"log"
	"os"
)

func main() {
	var server = flag.Bool("serve", false, "Run server")
	flag.Parse()
	env, err := initializers.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}
	logFile := initLog()
	defer logFile.Close()
	if *server {
		cmd.Serve(env)
	} else {
		fmt.Println("Running CLI")
		//cmd.ExecuteCli()
	}
}

func initLog() *os.File {
	logFilePath := os.Getenv("LOG_FILE_PATH")
	if logFilePath == "" {
		logFilePath = "."
	}
	logFilePath += "/log.txt"
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	w := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(w)
	return logFile
}
