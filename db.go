package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Conn struct {
	db *sql.DB
}

func (c *Conn) Init(config Config) error {
	dbase, err := sql.Open("postgres", config.DB)
	if err != nil {
		return err
	}

	c.db = dbase
	return c.db.Ping()
}
