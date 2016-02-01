// Reddit audiences crawler
// Rémy Mathieu © 2016
package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type TodayHandler struct {
	app *App
}

func (c TodayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	println(vars["subreddit"])
}
