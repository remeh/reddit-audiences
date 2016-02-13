// Reddit audiences crawler
// Rémy Mathieu © 2016
package web

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/remeh/reddit-audiences/app"
	"github.com/remeh/reddit-audiences/db"

	"github.com/pborman/uuid"
)

type SignupGet struct {
	App *app.App
}

type SignupPost struct {
	App *app.App
}

type signupParams struct {
	TemplateParams
	Email string
	Error string
}

func (c SignupGet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := c.App.Templates.Lookup("signup.html")
	t.Execute(w, signupParams{})
}

func (c SignupPost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := c.App.Templates.Lookup("signup.html")
	t_end := c.App.Templates.Lookup("signup_end.html")

	// read parameters
	// ----------------------

	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	passwordconfirm := r.Form.Get("passwordconfirm")

	// check parameters
	// ----------------------

	if len(email) == 0 ||
		!strings.Contains(email, ".") ||
		!strings.Contains(email, "@") {
		w.WriteHeader(400)
		t.Execute(w, signupParams{
			Email: email,
			Error: "Please fill a valid email.",
		})
		return
	}

	if len(password) == 0 {
		w.WriteHeader(400)
		t.Execute(w, signupParams{
			Email: email,
			Error: "Please fill a password.",
		})
		return
	}

	if len(passwordconfirm) == 0 {
		w.WriteHeader(400)
		t.Execute(w, signupParams{
			Email: email,
			Error: "Please confirm your password.",
		})
		return
	}

	if password != passwordconfirm {
		w.WriteHeader(400)
		t.Execute(w, signupParams{
			Email: email,
			Error: "Password confirmation doesn't match.",
		})
		return
	}

	if !app.IsPasswordSecure(password) {
		w.WriteHeader(400)
		t.Execute(w, signupParams{
			Email: email,
			Error: "The given password isn't strong enough.",
		})
		return
	}

	if exists, err := c.App.DB().ExistingEmail(email); err != nil {
		w.WriteHeader(500)
		t.Execute(w, signupParams{
			Email: email,
			Error: "An error occurred.",
		})
		log.Println("err: while crypting a password:", err.Error())
	} else if exists {
		w.WriteHeader(400)
		t.Execute(w, signupParams{
			Email: email,
			Error: "Existing email.",
		})
		return
	}

	// crypt the password
	// ----------------------

	cryptedPassword, err := app.CryptPassword(password)
	if err != nil {
		w.WriteHeader(500)
		t.Execute(w, signupParams{
			Email: email,
			Error: "An error occurred.",
		})
		log.Println("err: while crypting a password:", err.Error())
		return
	}

	// store the new user
	// ----------------------

	now := time.Now()
	user := db.User{
		Uuid:         uuid.New(),
		Email:        email,
		CreationTime: now,
		LastLogin:    now,
	}

	_, err = c.App.DB().InsertUser(user, cryptedPassword)
	if err != nil {
		w.WriteHeader(500)
		t.Execute(w, signupParams{
			Email: email,
			Error: "An error occurred.",
		})
		log.Printf("err: while creating an account for email '%s': %s", email, err.Error())
		return
	}

	// create the session and send the cookies.
	// ----------------------
	session, err := app.CreateSession(c.App.DB(), user, now)
	if err != nil {
		w.WriteHeader(500)
		t.Execute(w, signupParams{
			Email: email,
			Error: "An error occurred.",
		})
		log.Printf("err: while creating a session for email '%s': %s", email, err.Error())
		return
	}

	// set cookie
	app.SetSessionCookie(w, session)

	p := signupParams{
		TemplateParams: TemplateParams{User: app.User{Email: email}},
	}

	t_end.Execute(w, p)
}
