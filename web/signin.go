// Reddit audiences crawler
// Rémy Mathieu © 2016
package web

import (
	"log"
	"net/http"
	"time"

	"github.com/remeh/reddit-audiences/app"
	"github.com/remeh/reddit-audiences/db"

	"golang.org/x/crypto/bcrypt"
)

type SigninGet struct {
	App *app.App
}

type SigninPost struct {
	App *app.App
}

type signinParams struct {
	app.Params
	Email string
	Error string
}

func (c SigninGet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := c.App.Templates.Lookup("signin.html")
	t.Execute(w, signinParams{})
}

func (c SigninPost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := c.App.Templates.Lookup("signin.html")

	// read parameters
	// ----------------------

	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// check parameters
	// ----------------------
	ok, user, err := c.checkUserPassword(c.App, email, password)
	if err != nil || !ok {
		w.WriteHeader(403)
		t.Execute(w, signinParams{
			Email: email,
			Error: "Wrong password.",
		})
		return
	}

	// create the session and send the cookies.
	// ----------------------
	session, err := app.CreateSession(c.App.DB(), user, time.Now())
	if err != nil {
		w.WriteHeader(500)
		t.Execute(w, signinParams{
			Email: email,
			Error: "An error occurred.",
		})
		log.Printf("err: while creating a session for email '%s': %s", email, err.Error())
		return
	}

	// set cookie
	app.SetSessionCookie(w, session)

	http.Redirect(w, r, "/", 302)
}

func (c SigninPost) checkUserPassword(a *app.App, email, password string) (bool, db.User, error) {
	u, hash, err := a.DB().GetUserByEmail(email)
	if err != nil {
		return false, db.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, db.User{}, err
	}
	return err == nil, u, err
}
