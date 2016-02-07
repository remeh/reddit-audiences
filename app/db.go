// Reddit audiences crawler
// Rémy Mathieu © 2016
package app

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Conn struct {
	db *sql.DB
}

const (
	SUBREDDITS_TO_CRAWL = `
		SELECT "name"
		FROM "subreddit"
		WHERE 
			"last_crawl" <= $1
			AND
			"active" = true
		ORDER BY "last_crawl"
	`
	INSERT_SUBREDDIT_AUDIENCE = `
		INSERT INTO "audience"
		(subreddit, crawl_time, audience, subscribers)
		VALUES
		($1, $2, $3, $4)
	`
	INSERT_ARTICLE = `
		INSERT INTO "article"
		("subreddit", "article_id", "article_title", "article_link", "author", "rank", "crawl_time", "promoted", "sticky")
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	LAST_ARTICLE_STATE = `
		SELECT "article_id", "rank"
		FROM "article"
		WHERE
			"subreddit" = $1
			AND
			"article_id" = $2
		ORDER BY crawl_time DESC
		LIMIT 1
	`
	UPDATE_LAST_CRAWL_TIME = `
		UPDATE "subreddit"
		SET
			"last_crawl" = $2
		WHERE "name" = $1
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
	FIND_ARTICLES = `
		SELECT * FROM (
			SELECT DISTINCT ON ("subreddit", "article_id") "subreddit", "article_id", "article_title", "article_link", "author", "rank", "crawl_time", "promoted", "sticky"
			FROM
				"article"
			WHERE
				"subreddit" = $1
				AND
				"crawl_time" >= $2
				AND
				"crawl_time" <= $3
		) sub_query
		ORDER BY sub_query.rank
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

func (c *Conn) Init(config Config) error {
	dbase, err := sql.Open("postgres", config.DB)
	if err != nil {
		return err
	}

	c.db = dbase
	return c.db.Ping()
}

// FindArticleLastState returns for the given subreddit and
// article ID the article ID (if found) and its rank.
// If the returned article ID is empty, it means that the
// article hasn't been inserted in the DB yet.
func (c Conn) FindArticleLastState(subreddit, articleId string) (string, int, error) {
	var id string
	var rank int

	r, err := c.db.Query(LAST_ARTICLE_STATE, subreddit, articleId)
	if err != nil {
		return "", 0, err
	}

	if r == nil {
		return "", 0, nil
	}

	defer r.Close()

	if r.Next() {
		if err := r.Scan(&id, &rank); err != nil {
			return "", 0, err
		}
	}

	return id, rank, nil
}

func (c Conn) InsertArticle(article Article) (sql.Result, error) {
	return c.db.Exec(INSERT_ARTICLE, article.Subreddit, article.ArticleId, article.ArticleTitle, article.ArticleLink, article.Author, article.Rank, article.CrawlTime, article.Promoted, article.Sticky)
}

// GetSubredditsToCrawl returns the subreddits which must be
// crawled as soon as possible.
func (c Conn) FindSubredditsToCrawl(after time.Time) ([]string, error) {
	rv := make([]string, 0)

	r, err := c.db.Query(SUBREDDITS_TO_CRAWL, after)
	if err != nil {
		return rv, err
	}

	if r == nil {
		return rv, nil
	}

	defer r.Close()

	for r.Next() {
		var subreddit string
		if err := r.Scan(&subreddit); err != nil {
			return rv, err
		}

		if len(subreddit) > 0 {
			rv = append(rv, subreddit)
		}
	}

	return rv, nil
}

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
		var subreddit, articleId, articleTitle, articleLink, author string
		var rank int
		var crawlTime time.Time
		var promoted, sticky bool

		if err := r.Scan(&subreddit, &articleId, &articleTitle, &articleLink, &author, &rank, &crawlTime, &promoted, &sticky); err != nil {
			return rv, err
		}

		if len(articleId) > 0 {
			rv = append(rv, Article{
				Subreddit:    subreddit,
				ArticleId:    articleId,
				ArticleTitle: articleTitle,
				ArticleLink:  articleLink,
				Author:       author,
				Rank:         rank,
				CrawlTime:    crawlTime,
				Promoted:     promoted,
				Sticky:       sticky,
			})
		}
	}

	return rv, nil
}

func (c Conn) FindArticlesRanking(subreddit string, start, end time.Time) ([]Ranking, error) {
	rv := make([]Ranking, 0)

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

		if len(subreddit) > 0 {
			rv = append(rv, Ranking{
				Subreddit: subreddit,
				CrawlTime: crawlTime,
				Rank:      rank,
				ArticleId: articleId,
			})
		}
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

// InsertAudienceValue writes an audience value for the given subreddit
// and updates the last crawl time of the subreddit.
func (c Conn) InsertAudienceValue(subreddit string, audience, subscribers int64) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	now := time.Now()

	// write the value

	_, err = tx.Exec(INSERT_SUBREDDIT_AUDIENCE, subreddit, now, audience, subscribers)
	if err != nil {
		return err
	}

	// last crawl time

	_, err = tx.Exec(UPDATE_LAST_CRAWL_TIME, subreddit, now)
	if err != nil {
		return err
	}

	return tx.Commit()
}
