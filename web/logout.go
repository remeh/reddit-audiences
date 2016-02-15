// Reddit audiences crawler
// Rémy Mathieu © 2016
package web

import (
	"net/http"
	"strings"

	"github.com/remeh/reddit-audiences/app"
)

type Logout struct {
	App *app.App
}

func (c Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// delete session
	// ----------------------
	var t *http.Cookie
	var err error
	var v string
	if t, err = r.Cookie("t"); err == nil {
		v = strings.Trim(t.Value, " ")
		c.App.DB().DeleteSession(v)
	}

	// delete cookie
	// ----------------------
	cookie := &http.Cookie{
		Name:   "t",
		Value:  v,
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 302)
}
