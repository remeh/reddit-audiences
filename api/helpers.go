// Reddit audiences crawler
// Rémy Mathieu © 2016
package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func render(w http.ResponseWriter, code int, r interface{}) {
	if d, err := json.Marshal(r); err != nil {
		w.WriteHeader(500)
		log.Printf("err: while rendering: %s", err.Error())
		return
	} else {
		w.WriteHeader(code)
		w.Write(d)
		return
	}
}
