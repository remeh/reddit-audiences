package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	REDDIT_SUBREDDIT_URL = "https://reddit.com/r/"
)

func StartCrawlingJob(a *App) {
	log.Println("info: starts tracking job.")
	ticker := time.NewTicker(time.Minute * 1)
	for range ticker.C {
		log.Println("info: tracking job is running.")
		Crawl(a)
	}
	ticker.Stop()
}

func Crawl(a *App) {
	// crawl each subreddit each 5 minutes
	five := time.Minute * 5
	t := time.Now().Add(-five)
	subreddits, err := a.DB().FindSubredditsToCrawl(t)

	if err != nil {
		log.Printf("err: can't retrieve subreddits to crawl: %s\n", err.Error())
	}

	for _, subreddit := range subreddits {
		log.Println("Crawling", subreddit)
		go func() {
			if audience, err := GetAudience(subreddit); err == nil {
				if err := a.DB().InsertSubredditValue(subreddit, audience); err != nil {
					log.Println("err:", err.Error())
				} else {
					log.Printf("info: subreddit %s has %d active users\n", subreddit, audience)
				}
			} else if err != nil {
				log.Println("err:", err.Error())
			}
		}()
	}
}

func GetAudience(subreddit string) (int, error) {
	var audience int
	var err error

	doc, err := goquery.NewDocument(REDDIT_SUBREDDIT_URL + subreddit)
	if err != nil {
		return 0, err
	}

	doc.Find("p.users-online span.number").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			log.Println("warn: found many times the number value.")
			return
		}

		// it looks like we found a value in the dom
		value := s.Text()
		if len(value) == 0 {
			err = fmt.Errorf("can't retrieve subreddit %s audience", subreddit)
			return
		}

		// sometimes it starts with ~
		if strings.HasPrefix(value, "~") {
			value = value[1:]
		}
		// , for thousands etc.
		value = strings.Replace(value, ",", "", -1)
		// finally trim
		value = strings.Trim(value, " ")

		audience, err = strconv.Atoi(value)
	})

	return audience, err
}
