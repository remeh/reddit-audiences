// Reddit audiences crawler
// Rémy Mathieu © 2016
package app

import (
	"log"
	"time"
)

const (
	SESSION_EXPIRATION = "2m"
)

func StartCleanSessionsJob(a *App) {
	duration, err := time.ParseDuration(SESSION_EXPIRATION)
	if err != nil {
		log.Println("err: the SESSION_EXPIRATION const is malformed:", SESSION_EXPIRATION)
		log.Println("err: the sessions won't expire.")
		return
	}

	ticker := time.NewTicker(time.Minute)

	log.Println("info: clean sessions job started.")

	for range ticker.C {
		cleanSessions(a, duration)
	}
	ticker.Stop()
}

func cleanSessions(a *App, sessionExpiration time.Duration) {
	if result, err := a.DB().DeleteExpiredSessions(sessionExpiration); err == nil {
		if ra, err := result.RowsAffected(); err == nil {
			if ra > 0 {
				log.Printf("info: %d session expired.\n", ra)
			}
		}
	}
}
