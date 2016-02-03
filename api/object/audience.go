package object

import (
	"time"

	"github.com/remeh/reddit-audiences/app"
)

type Audience struct {
	CrawlTime time.Time `json:"crawl_time"`
	Audience  int64     `json:"audience"`
}

func AudiencesFromApp(audiences []app.Audience) []Audience {
	rv := make([]Audience, len(audiences))
	for i, a := range audiences {
		rv[i] = AudienceFromApp(a)
	}
	return rv
}

func AudienceFromApp(audience app.Audience) Audience {
	return Audience{
		CrawlTime: audience.CrawlTime,
		Audience:  audience.Audience,
	}
}
