// Reddit audiences crawler
// Rémy Mathieu © 2016
package db

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
		("subreddit", "article_id", "article_title", "article_link", "article_external_link", "author", "rank", "crawl_time", "promoted", "sticky")
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
)

func (c *Conn) Init(connString string) error {
	dbase, err := sql.Open("postgres", connString)
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
	return c.db.Exec(INSERT_ARTICLE, article.Subreddit, article.ArticleId, article.ArticleTitle, article.ArticleLink, article.ArticleExternalLink, article.Author, article.Rank, article.CrawlTime, article.Promoted, article.Sticky)
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
