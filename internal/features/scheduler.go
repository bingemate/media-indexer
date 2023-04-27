package features

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func ScheduleScanner(cronStr string, movieScanner *MovieScanner, tvScanner *TVScanner) {
	cronTab, err := cron.ParseStandard(cronStr)
	if err != nil {
		log.Println("Disabling scanner scheduling due to invalid cron expression:", err)
		return
	}
	c := cron.New()
	_, err = c.AddFunc(cronStr, func() {
		log.Println("Scanning for new media...")
		movies, err := movieScanner.ScanMovies()
		if err != nil {
			log.Println("Error scanning movies:", err)
		}
		log.Println("Found", len(*movies), "new movies")
		tvs, err := tvScanner.ScanTV()
		if err != nil {
			log.Println("Error scanning tvs:", err)
		}
		log.Println("Found", len(*tvs), "new tvs")
		log.Println("Next scan scheduled for", cronTab.Next(time.Now()).Format(time.RFC1123))
	})
	if err != nil {
		return
	}
	c.Start()
	log.Println("Next scan scheduled for", cronTab.Next(time.Now()).Format(time.RFC1123))
	select {}
}
