// Reddit audiences crawler
// Rémy Mathieu © 2016
package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func render(w http.ResponseWriter, code int, r interface{}) {
	if d, err := json.Marshal(r); err != nil {
		w.WriteHeader(500)
		log.Printf("err: while rendering: %s", err.Error())
		return
	} else {
		w.Write(d)
		w.WriteHeader(200)
		return
	}
}
