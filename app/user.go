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

type User struct {
	Email     string
	Firstname string
	Lastname  string
}

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
	w.Header().Set("Set-Cookie", fmt.Sprintf("t=%s", session.Token))
}

// ----------------------

type TemplateParams struct {
	User User
}

func TmplParams(app *App, r *http.Request) TemplateParams {
	return TemplateParams{
		User: GetUser(app.DB(), r),
	}
}

// ----------------------

func GetUser(conn db.Conn, r *http.Request) User {
	if r == nil {
		return User{}
	}

	cookie, err := r.Cookie("t")
	if err != nil {
		return User{}
	}

	sessionToken := cookie.Value

	user, err := conn.GetUserFromSessionToken(sessionToken)
	if err != nil {
		log.Printf("err: while getting an user from the session ID '%s': %s", sessionToken, err.Error())
		return User{}
	}

	conn.UpdateSession(sessionToken)

	return User{
		Email:     user.Email,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
	}
}
