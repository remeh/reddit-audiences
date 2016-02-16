// Reddit audiences crawler
// Rémy Mathieu © 2016
package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/remeh/reddit-audiences/app"

	"github.com/gorilla/mux"
)

type AnnotateHandler struct {
	App *app.App
}

type annotateBody struct {
	Time    time.Time `json:"t"`
	Message string    `json:"m"`
}

func (c AnnotateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	subreddit := vars["subreddit"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("err: while reading body of annotate: %s\n", err.Error())
		return
	}

	defer r.Body.Close()

	var body annotateBody
	if err := json.Unmarshal(data, &body); err != nil {
		w.WriteHeader(400)
		return
	}

	if len(subreddit) == 0 || body.Time.IsZero() || len(body.Message) == 0 {
		w.WriteHeader(400)
		return
	}

	u := app.GetUser(c.App.DB(), r)
	if len(u.Email) == 0 {
		w.WriteHeader(403)
		return
	}

	err = c.App.DB().InsertAnnotation(u.Uuid, subreddit, body.Time, body.Message)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("err: while annotating by '%s' on '%s': %s", u.Email, subreddit, err.Error())
		return
	}
}
