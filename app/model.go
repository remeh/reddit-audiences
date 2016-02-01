package app

import (
	"time"
)

type Audience struct {
	Subreddit string
	Audience  int
	CrawlTime time.Time
}
