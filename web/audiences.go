// Reddit audiences crawler
// Rémy Mathieu © 2016
package web

import (
	"net/http"
	"strings"

	"github.com/remeh/reddit-audiences/app"

	"github.com/gorilla/mux"
)

type Audiences struct {
	App *app.App
}

type audiencesParams struct {
	app.Params
	Subreddit string
}

func (c Audiences) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := c.App.Templates.Lookup("audiences.html")

	vars := mux.Vars(r)
	subreddit := vars["subreddit"]

	// test if known subreddit
	subreddit = strings.ToLower(strings.Trim(subreddit, " "))
	if len(subreddit) == 0 {
		http.Redirect(w, r, "/", 301)
		return
	}

	t.Execute(w, audiencesParams{
		Params:    app.TmplParams(c.App, r, "Audiences"),
		Subreddit: subreddit,
	})
}
