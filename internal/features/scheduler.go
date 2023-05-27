package features

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

import "sync"

var mutex = &sync.Mutex{}

func ScheduleScanner(cronStr string, movieScanner *MovieScanner, tvScanner *TVScanner) {
	cronTab, err := cron.ParseStandard(cronStr)
	if err != nil {
		log.Println("Disabling scanner scheduling due to invalid cron expression:", err)
		return
	}
	c := cron.New()
	_, err = c.AddFunc(cronStr, func() {
		// Verrouillage de la mutex
		locked := mutex.TryLock()
		if !locked {
			log.Println("Scanner already running, skipping this run")
			return
		}
		defer mutex.Unlock()

		log.Println("Scanning for new media...")
		movies, err := movieScanner.ScanMovies()
		if err != nil {
			log.Println("Error scanning movies:", err)
		} else {
			log.Println("Found", len(*movies), "new movies")
		}
		tvs, err := tvScanner.ScanTV()
		if err != nil {
			log.Println("Error scanning tvs:", err)
		} else {
			log.Println("Found", len(*tvs), "new tvs")
		}
		log.Println("Next scan scheduled for", cronTab.Next(time.Now()).Format(time.RFC1123))
	})
	if err != nil {
		return
	}
	c.Start()
	log.Println("Next scan scheduled for", cronTab.Next(time.Now()).Format(time.RFC1123))
}
