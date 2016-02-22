package db

import (
	"time"
)

type ArticleState string

var (
	Rising   ArticleState = "rising"
	Stagnant ArticleState = "stagnant"
	Falling  ArticleState = "falling"
	New      ArticleState = "new"
	Removed  ArticleState = "removed"
)

type Audience struct {
	Subreddit   string
	Audience    int64
	Subscribers int64
	CrawlTime   time.Time
}

type Article struct {
	Subreddit           string
	ArticleId           string
	ArticleTitle        string
	ArticleLink         string
	ArticleExternalLink string
	Score               int
	Comments            int
	Author              string
	Rank                int
	CrawlTime           time.Time
	Promoted            bool
	Sticky              bool
}

type Annotation struct {
	Owner     string
	Subreddit string
	Message   string
	Time      time.Time
}

type Ranking struct {
	Subreddit string
	ArticleId string
	CrawlTime time.Time
	Rank      int
	Score     int
	Comments  int
}

type User struct {
	Uuid         string
	Email        string
	Firstname    string
	Lastname     string
	CreationTime time.Time
	LastLogin    time.Time
}

type Session struct {
	Token   string
	User    User
	HitTime time.Time
}

// ----------------------

type Rankings []Ranking

func (r Rankings) Len() int           { return len(r) }
func (r Rankings) Swap(a, b int)      { r[a], r[b] = r[b], r[a] }
func (r Rankings) Less(a, b int) bool { return r[a].CrawlTime.Before(r[b].CrawlTime) }
