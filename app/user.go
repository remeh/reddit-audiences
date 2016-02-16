// Reddit audiences crawler
// Rémy Mathieu © 2016
package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/remeh/reddit-audiences/db"

	"github.com/pborman/uuid"
)

// CreationSessions creates in-base a session for the
// given user already created in database.
func CreateSession(conn db.Conn, user db.User, creationTime time.Time) (db.Session, error) {
	if len(user.Uuid) == 0 {
		return db.Session{}, fmt.Errorf("nil user given to CreateSession")
	}

	session := db.Session{
		Token:   uuid.New(),
		User:    user,
		HitTime: creationTime,
	}

	err := conn.InsertSession(session)
	return session, err
}

func SetSessionCookie(w http.ResponseWriter, session db.Session) {
	cookie := &http.Cookie{
		Name:   "t",
		Value:  session.Token,
		MaxAge: 86400, // 1 day
	}
	http.SetCookie(w, cookie)
}

// ----------------------

func GetUser(conn db.Conn, r *http.Request) db.User {
	if r == nil {
		return db.User{}
	}

	cookie, err := r.Cookie("t")
	if err != nil {
		return db.User{}
	}

	sessionToken := cookie.Value

	user, err := conn.GetUserFromSessionToken(sessionToken)
	if err != nil {
		log.Printf("err: while getting an user from the session ID '%s': %s", sessionToken, err.Error())
		return db.User{}
	}

	conn.UpdateSession(sessionToken)

	return user
}
