package web

import (
	"net/http"

	"github.com/remeh/reddit-audiences/app"
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

	t = t.Funcs(app.TemplateHelpers())
	t.Execute(w, signup{
		Email: "mail@mail.com", // TODO(remy),
		Error: "Please fill a valid email.",
	})
}
