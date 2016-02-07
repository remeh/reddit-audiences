package app

import (
	"time"
)

type Audience struct {
	Subreddit   string
	Audience    int64
	Subscribers int64
	CrawlTime   time.Time
}

type Article struct {
	Subreddit    string
	ArticleId    string
	ArticleTitle string
	ArticleLink  string
	Author       string
	Rank         int
	CrawlTime    time.Time
	Promoted     bool
	Sticky       bool
}

type Ranking struct {
	Subreddit string
	ArticleId string
	CrawlTime time.Time
	Rank      int
}
