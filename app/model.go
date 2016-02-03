package app

import (
	"time"
)

type Audience struct {
	Subreddit string
	Audience  int64
	CrawlTime time.Time
}
