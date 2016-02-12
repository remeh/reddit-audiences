package object

import (
	"time"

	"github.com/remeh/reddit-audiences/db"
)

type Rankings map[string][]db.Ranking

type Ranking struct {
	CrawlTime time.Time `json:"crawl_time"`
	ArticleId string    `json:"article_id"`
	Rank      int       `json:"rank"`
}

func RankingsFromApp(rankings []db.Ranking) Rankings {
	rv := make(Rankings)

	for _, r := range rankings {
		var exists bool

		if _, exists = rv[r.ArticleId]; !exists {
			rv[r.ArticleId] = make([]db.Ranking, 0)
		}

		rv[r.ArticleId] = append(rv[r.ArticleId], RankingFromApp(r))
	}

	return rv
}

func RankingFromApp(ranking db.Ranking) db.Ranking {
	return db.Ranking{
		CrawlTime: ranking.CrawlTime,
		Rank:      ranking.Rank,
		ArticleId: ranking.ArticleId,
	}
}
