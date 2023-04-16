package initializers

import (
	"io"
	"log"
	"os"
)

func InitLog(logfile string) *os.File {
	logFile, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	w := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(w)
	return logFile
}
