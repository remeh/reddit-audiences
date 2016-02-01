// Reddit audiences crawler
// Rémy Mathieu © 2016
package main

import (
	"log"
	"net/http"
)

type LogAdapter struct {
	app     *App
	handler http.Handler
}

func (a LogAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// then propagate to the next handler if no 403 has been returned.
	sWriter := &StatusWriter{w, 200}
	a.handler.ServeHTTP(sWriter, r)

	w.Header().Set("Content-Type", "application/json")

	log.Printf("info: hit: %s %s %s referer[%s] user-agent[%s] addr[%s] code[%d]\n", r.Method, r.URL.String(), r.Proto, r.Referer(), r.UserAgent(), r.RemoteAddr, sWriter.Status)
}

// LogRoute creates a route which will log the route access.
func LogRoute(a *App, handler http.Handler) http.Handler {
	return LogAdapter{
		app:     a,
		handler: handler,
	}
}

type StatusWriter struct {
	http.ResponseWriter
	Status int
}

func (w *StatusWriter) WriteHeader(code int) {
	w.Status = code
	w.ResponseWriter.WriteHeader(code)
}
