// Reddit audiences crawler
// Rémy Mathieu © 2016
package db

import (
	"database/sql"
	"time"
)

const (
	INSERT_USER = `
		INSERT INTO "user"
		("uuid", "email", "hash", "firstname", "lastname", "creation_time", "last_login")
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
	`

	INSERT_SESSION = `
		INSERT INTO "session"
		("token", "uuid", "creation_time")
		VALUES
		($1, $2, $3)
	`

	EXISTING_EMAIL = `
		SELECT "email"
		FROM "user"
		WHERE "email" = $1
	`

	USER_FROM_SESSION_UUID = `
		SELECT "user"."uuid", "user"."email", "user"."firstname", "user"."lastname", "user"."creation_time", "user"."last_login"
		FROM "user"
		JOIN "session" ON "session"."uuid" = "user"."uuid" AND "session"."token" = $1
	`
)

func (c Conn) InsertUser(user User, hash string) (sql.Result, error) {
	return c.db.Exec(INSERT_USER, user.Uuid, user.Email, hash, user.Firstname, user.Lastname, user.CreationTime, user.LastLogin)
}

func (c Conn) InsertSession(session Session) error {
	_, err := c.db.Exec(INSERT_SESSION, session.Token, session.User.Uuid, session.CreationTime)
	return err
}

func (c Conn) ExistingEmail(email string) (bool, error) {
	r, err := c.db.Query(EXISTING_EMAIL, email)
	if err != nil {
		return false, err
	}

	if r == nil {
		return false, nil
	}

	defer r.Close()

	if r.Next() {
		return true, nil
	}
	return false, nil
}

func (c Conn) GetUserFromSessionToken(token string) (User, error) {
	r, err := c.db.Query(USER_FROM_SESSION_UUID, token)
	if err != nil {
		return User{}, err
	}

	if r == nil {
		return User{}, nil
	}

	defer r.Close()

	if r.Next() {
		var uuid, email, firstname, lastname string
		var creationTime, lastLogin time.Time

		if err := r.Scan(&uuid, &email, &firstname, &lastname, &creationTime, &lastLogin); err == nil {
			return User{
				Uuid:         uuid,
				Email:        email,
				Firstname:    firstname,
				Lastname:     lastname,
				CreationTime: creationTime,
				LastLogin:    lastLogin,
			}, nil
		} else {
			return User{}, err
		}
	}

	return User{}, nil
}
