package initializers

import (
	"io"
	"log"
	"os"
)

func InitLog() *os.File {
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
