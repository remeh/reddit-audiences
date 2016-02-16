package object

import (
	"time"

	"github.com/remeh/reddit-audiences/db"
)

type Audience struct {
	CrawlTime time.Time `json:"crawl_time"`
	Audience  int64     `json:"audience"`
}

func AudiencesFromApp(audiences []db.Audience) []Audience {
	rv := make([]Audience, len(audiences))
	for i, a := range audiences {
		rv[i] = AudienceFromApp(a)
	}
	return rv
}

func AudienceFromApp(audience db.Audience) Audience {
	return Audience{
		CrawlTime: audience.CrawlTime,
		Audience:  audience.Audience,
	}
}

type Annotation struct {
	Time    time.Time `json:"crawl_time"`
	Message string    `json:"message"`
}
