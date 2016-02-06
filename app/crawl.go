// Reddit audiences crawler
// Rémy Mathieu © 2016
package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	REDDIT_SUBREDDIT_URL = "https://reddit.com/r/"
)

func StartCrawlingJob(a *App) {
	if a.Config.Crawl {
		log.Println("info: starts tracking job.")
		ticker := time.NewTicker(time.Second * 30)
		for range ticker.C {
			log.Println("info: tracking job is running.")
			Crawl(a)
		}
		ticker.Stop()
	}
}

// Crawl retrieves the audience of subreddits for which
// the last crawl time is more than some minutes.
func Crawl(a *App) {
	// crawl each subreddit each 5 minutes
	five := time.Minute * 5
	t := time.Now().Add(-five)
	subreddits, err := a.DB().FindSubredditsToCrawl(t)

	if err != nil {
		log.Printf("err: can't retrieve subreddits to crawl: %s\n", err.Error())
	}

	for _, subreddit := range subreddits {
		log.Println("info: crawling", subreddit)
		go func(subreddit string) {
			if audience, subscribers, err := GetAudience(subreddit); err == nil {
				// store the value and update the last crawl time
				if err := a.DB().InsertAudienceValue(subreddit, audience, subscribers); err != nil {
					log.Println("err:", err.Error())
				} else {
					log.Printf("info: subreddit %s has %d active users (%d subscribers)\n", subreddit, audience, subscribers)
				}
			} else if err != nil {
				log.Println("err:", err.Error())
			}
		}(subreddit)
	}
}

// GetAudience gets the subreddit page on reddit
// and gets the current audience of this subreddit in the DOM.
// NOTE(remy): we stop as soon as we have a DOM error because
// it has great chances that the full DOM is corrupted/not retrieved.
func GetAudience(subreddit string) (int64, int64, error) {
	var audience int64
	var subscribers int64
	var err error

	doc, err := getSubredditPage(REDDIT_SUBREDDIT_URL + subreddit)
	if err != nil {
		return 0, 0, fmt.Errorf("while crawling %s: %s", subreddit, err.Error())
	}

	// audience
	// ----------------------

	s := doc.Find("p.users-online span.number").First()

	value := s.Text()
	if len(value) != 0 {
		return 0, 0, fmt.Errorf("can't retrieve subreddit %s audience: no text value in the dom node.", subreddit)
	}

	if audience, err = cleanInt(value); err != nil {
		return 0, 0, err
	}

	// subscribers
	// ----------------------

	s = doc.Find("span.subscribers span.number").First()

	value = s.Text()
	if len(value) == 0 {
		return 0, 0, fmt.Errorf("can't retrieve subreddit %s subscribers: no text value in the dom node.", subreddit)
	}

	if subscribers, err = cleanInt(value); err != nil {
		return 0, 0, err
	}

	return audience, subscribers, err
}

func getSubredditPage(url string) (*goquery.Document, error) {
	// TODO(remy): create a pool of clients and use it
	client := http.Client{}

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// TODO(remy)
	r.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:42.0) Gecko/20100101 Firefox/42.0")

	resp, err := client.Do(r)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(resp.Status)
	}

	return goquery.NewDocumentFromResponse(resp)
}

func cleanInt(str string) (int64, error) {
	// sometimes it starts with ~
	if strings.HasPrefix(str, "~") {
		str = str[1:]
	}
	// , for thousands etc.
	str = strings.Replace(str, ",", "", -1)
	// finally trim
	str = strings.Trim(str, " ")

	return strconv.ParseInt(str, 10, 64)
}
