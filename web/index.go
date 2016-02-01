package web

import (
	"net/http"

	"github.com/remeh/reddit-audiences/app"

	"github.com/gorilla/mux"
)

type Index struct {
	App *app.App
}

func (c Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := c.App.Templates.Lookup("index.html")

	vars := mux.Vars(r)
	subreddit := vars["subreddit"]

	if len(subreddit) == 0 {
		// TODO(remy): redirect to somewhere to enter the subreddit
	}

	t.Execute(w, nil)
}
