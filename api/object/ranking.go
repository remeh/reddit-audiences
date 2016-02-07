package object

import (
	"time"

	"github.com/remeh/reddit-audiences/app"
)

type Rankings map[string][]Ranking

type Ranking struct {
	CrawlTime time.Time `json:"crawl_time"`
	ArticleId string    `json:"article_id"`
	Rank      int       `json:"rank"`
}

func RankingsFromApp(rankings []app.Ranking) Rankings {
	rv := make(Rankings)

	for _, r := range rankings {
		var exists bool

		if _, exists = rv[r.ArticleId]; !exists {
			rv[r.ArticleId] = make([]Ranking, 0)
		}

		rv[r.ArticleId] = append(rv[r.ArticleId], RankingFromApp(r))
	}

	return rv
}

func RankingFromApp(ranking app.Ranking) Ranking {
	return Ranking{
		CrawlTime: ranking.CrawlTime,
		Rank:      ranking.Rank,
		ArticleId: ranking.ArticleId,
	}
}
