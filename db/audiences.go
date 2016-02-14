// Reddit audiences crawler
// Rémy Mathieu © 2016
package db

import (
	"time"

	_ "github.com/lib/pq"
)

const (
	FIND_ARTICLES = `
		SELECT * FROM (
			SELECT DISTINCT ON ("subreddit", "article_id") "subreddit", "article_id", "article_title", "article_link", "article_external_link", "author", "rank", "crawl_time", "promoted", "sticky"
			FROM
				"article"
			WHERE
				"subreddit" = $1
				AND
				"crawl_time" >= $2
				AND
				"crawl_time" <= $3
		) sub_query
		ORDER BY sub_query.rank, crawl_time DESC
	`
	ARTICLES_RANKING = `
		SELECT "rank", "article_id", "crawl_time"
		FROM "article"
		WHERE
			"subreddit" = $1
			AND
			"crawl_time" >= $2
			AND
			"crawl_time" <= $3
		ORDER BY "crawl_time"
	`
	AUDIENCES_INTERVAL = `
		SELECT "audience", "crawl_time"
		FROM "audience"
		WHERE
			"subreddit" = $1
			AND
			"crawl_time" >= $2
			AND
			"crawl_time" <= $3
		ORDER BY "crawl_time"
	`
)

func (c Conn) FindArticles(subreddit string, start, end time.Time) ([]Article, error) {
	rv := make([]Article, 0)

	r, err := c.db.Query(FIND_ARTICLES, subreddit, start, end)
	if err != nil {
		return rv, err
	}

	if r == nil {
		return rv, nil
	}

	defer r.Close()

	for r.Next() {
		var subreddit, articleId, articleTitle, articleLink, articleExtLink, author string
		var rank int
		var crawlTime time.Time
		var promoted, sticky bool

		if err := r.Scan(&subreddit, &articleId, &articleTitle, &articleLink, &articleExtLink, &author, &rank, &crawlTime, &promoted, &sticky); err != nil {
			return rv, err
		}

		if len(articleId) > 0 {
			rv = append(rv, Article{
				Subreddit:           subreddit,
				ArticleId:           articleId,
				ArticleTitle:        articleTitle,
				ArticleLink:         articleLink,
				ArticleExternalLink: articleExtLink,
				Author:              author,
				Rank:                rank,
				CrawlTime:           crawlTime,
				Promoted:            promoted,
				Sticky:              sticky,
			})
		}
	}

	return rv, nil
}

func (c Conn) FindArticlesRanking(subreddit string, start, end time.Time) (map[string][]Ranking, error) {
	rv := make(map[string][]Ranking)

	r, err := c.db.Query(ARTICLES_RANKING, subreddit, start, end)
	if err != nil {
		return rv, err
	}

	if r == nil {
		return rv, nil
	}

	defer r.Close()

	for r.Next() {
		var rank int
		var articleId string
		var crawlTime time.Time

		if err := r.Scan(&rank, &articleId, &crawlTime); err != nil {
			return rv, err
		}

		if len(articleId) == 0 {
			continue
		}

		if _, exists := rv[articleId]; !exists {
			rv[articleId] = make([]Ranking, 0)
		}

		rv[articleId] = append(rv[articleId], Ranking{
			Subreddit: subreddit,
			CrawlTime: crawlTime,
			Rank:      rank,
			ArticleId: articleId,
		})
	}

	return rv, nil
}

func (c Conn) FindAudiencesInterval(subreddit string, start, end time.Time) ([]Audience, error) {
	rv := make([]Audience, 0)

	r, err := c.db.Query(AUDIENCES_INTERVAL, subreddit, start, end)
	if err != nil {
		return rv, err
	}

	if r == nil {
		return rv, nil
	}

	defer r.Close()

	for r.Next() {
		var audience int64
		var crawlTime time.Time

		if err := r.Scan(&audience, &crawlTime); err != nil {
			return rv, err
		}

		if len(subreddit) > 0 {
			rv = append(rv, Audience{
				Subreddit: subreddit,
				CrawlTime: crawlTime,
				Audience:  audience,
			})
		}
	}

	return rv, nil
}
