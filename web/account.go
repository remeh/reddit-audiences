// Reddit audiences crawler
// Rémy Mathieu © 2016
package web

import (
	"net/http"

	"github.com/remeh/reddit-audiences/app"
)

type Account struct {
	App *app.App
}

type accountParams struct {
	app.Params
}

func (c Account) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := c.App.Templates.Lookup("account.html")
	p := app.TmplParams(c.App, r, "Account")

	if len(p.User.Email) == 0 {
		http.Redirect(w, r, "/signin", 302)
	}

	t.Execute(w, accountParams{
		p,
	})
}
