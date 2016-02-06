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
	UPDATE_LAST_CRAWL_TIME = `
		UPDATE "subreddit"
		SET
			"last_crawl" = $2
		WHERE "name" = $1
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
