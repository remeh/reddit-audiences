package web

import (
	"net/http"

	"github.com/remeh/reddit-audiences/app"

	"github.com/gorilla/mux"
)

type Audiences struct {
	App *app.App
}

type audiencesBody struct {
	Subreddit string
}

func (c Audiences) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := c.App.Templates.Lookup("audiences.html")

	vars := mux.Vars(r)
	subreddit := vars["subreddit"]

	// test if known subreddit

	if len(subreddit) == 0 {
		http.Redirect(w, r, "/", 301)
		return
	}

	t.Execute(w, audiencesBody{
		Subreddit: Capitalize(subreddit),
	})
}
