// Reddit audiences crawler
// Rémy Mathieu © 2016
package web

import (
	"log"
	"net/http"

	"github.com/remeh/reddit-audiences/app"
	"github.com/remeh/reddit-audiences/db"
)

type TemplateParams struct {
	User app.User
}

func templateParams(conn db.Conn, r *http.Request) TemplateParams {
	return TemplateParams{
		User: GetUser(conn, r),
	}
}

func GetUser(conn db.Conn, r *http.Request) app.User {
	if r == nil {
		return app.User{}
	}

	cookie, err := r.Cookie("t")
	if err != nil {
		return app.User{}
	}

	sessionToken := cookie.Value

	user, err := conn.GetUserFromSessionToken(sessionToken)
	if err != nil {
		log.Printf("err: while getting an user from the session ID '%s': %s", sessionToken, err.Error())
		return app.User{}
	}

	conn.UpdateSession(sessionToken)

	return app.User{
		Email:     user.Email,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
	}
}
