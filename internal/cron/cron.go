package cron

import (
	"log"
	"portfolyo/internal/service"
	"time"

	"github.com/robfig/cron/v3"
)

func Start(ks service.KurService) {
	loc, _ := time.LoadLocation("Europe/Istanbul")
	c := cron.New(cron.WithLocation(loc))

	_, err := c.AddFunc("0 17 * * *", func() { // "*/1 * * * *" --> her dakikada || "0 17 * * *" --> her gün saat 17 de
		log.Println("cron tetiklendi kral")
		err := ks.FetchAndSaveRates()
		if err != nil {
			log.Println("cron error:", err)
		} else {
			log.Println("cron: exchange rates updated successfully")
		}
	})

	if err != nil {
		log.Println("cron setup error:", err)
	}

	c.Start()
}
