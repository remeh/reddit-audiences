package main

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
			"next_crawl" <= $1
		ORDER BY "last_crawl"
	`
	INSERT_SUBREDDIT_AUDIENCE = `
		INSERT INTO "audience"
		(subreddit, crawl_time, audience)
		VALUES
		($1, $2, $3)
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
func (c *Conn) FindSubredditsToCrawl() ([]string, error) {
	rv := make([]string, 0)

	r, err := c.db.Query(SUBREDDITS_TO_CRAWL, time.Now())
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

// InsertSubredditValue writes an audience value for the given subreddit.
func (c *Conn) InsertSubredditValue(subreddit string, value uint) error {
	_, err := c.db.Exec(INSERT_SUBREDDIT_AUDIENCE, subreddit, time.Now(), value)
	if err != nil {
		return err
	}
	return nil
}
