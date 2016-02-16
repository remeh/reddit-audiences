// Reddit audiences crawler
// Rémy Mathieu © 2016
package api

import (
	"log"
	"net/http"
	"time"

	"github.com/remeh/reddit-audiences/app"

	"github.com/gorilla/mux"
)

type AnnotateHandler struct {
	App *app.App
}

func (c AnnotateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	subreddit := vars["subreddit"]

	r.ParseForm()
	t := r.Form.Get("t")
	message := r.Form.Get("m")

	if len(subreddit) == 0 || len(t) == 0 || len(message) == 0 {
		w.WriteHeader(400)
		return
	}

	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	u := app.GetUser(c.App.DB(), r)
	if len(u.Email) == 0 {
		w.WriteHeader(403)
		return
	}

	err = c.App.DB().InsertAnnotation(u.Uuid, subreddit, parsedTime, message)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("err: while annotating by '%s' on '%s': %s", u.Email, subreddit, message)
		return
	}
}
