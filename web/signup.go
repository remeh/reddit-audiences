package web

import (
	"log"
	"net/http"
	"strings"

	"github.com/remeh/reddit-audiences/app"
	"github.com/remeh/reddit-audiences/db"
)

type SignupGet struct {
	App *app.App
}

type SignupPost struct {
	App *app.App
}

type signup struct {
	Email string
	Error string
}

func (c SignupGet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := c.App.Templates.Lookup("signup.html")

	t = t.Funcs(app.TemplateHelpers())
	t.Execute(w, signup{})
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
		t.Execute(w, signup{
			Email: email,
			Error: "Please fill a valid email.",
		})
		return
	}

	if len(password) == 0 {
		w.WriteHeader(400)
		t.Execute(w, signup{
			Email: email,
			Error: "Please fill a password.",
		})
		return
	}

	if len(passwordconfirm) == 0 {
		w.WriteHeader(400)
		t.Execute(w, signup{
			Email: email,
			Error: "Please confirm your password.",
		})
		return
	}

	if password != passwordconfirm {
		w.WriteHeader(400)
		t.Execute(w, signup{
			Email: email,
			Error: "Password confirmation doesn't match.",
		})
		return
	}

	if !app.IsPasswordSecure(password) {
		w.WriteHeader(400)
		t.Execute(w, signup{
			Email: email,
			Error: "The given password isn't strong enough.",
		})
		return
	}

	// crypt the password
	// ----------------------

	cryptedPassword, err := app.CryptPassword(password)
	if err != nil {
		w.WriteHeader(500)
		t.Execute(w, signup{
			Email: email,
			Error: "An error occurred.",
		})
		log.Println("err: while crypting a password:", err.Error())
		return
	}

	// store the new user
	// ----------------------

	// TODO(remy): store and response
	user := db.User{
		Email: email,
	}

	_, err = c.App.DB().InsertUser(user, cryptedPassword)
	if err != nil {
		w.WriteHeader(500)
		t.Execute(w, signup{
			Email: email,
			Error: "An error occurred.",
		})
		log.Println("err: while creating an account for email:", email)
	}

	t_end.Execute(w, nil)
}
