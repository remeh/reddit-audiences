// Reddit audiences crawler
// Rémy Mathieu © 2016
package db

import (
	"database/sql"
)

const (
	INSERT_USER = `
		INSERT INTO "user"
		(email, hash, firstname, lastname)
		VALUES
		($1, $2, $3, $4)
	`
)

func (c Conn) InsertUser(user User, hash string) (sql.Result, error) {
	return c.db.Exec(INSERT_USER, user.Email, hash, user.Firstname, user.Lastname)
}
