// Reddit audiences crawler
// Rémy Mathieu © 2016
package web

import (
	"net/http"

	"github.com/remeh/reddit-audiences/app"
)

type Index struct {
	App *app.App
}

type indexParams struct {
	app.Params
}

func (c Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := c.App.Templates.Lookup("index.html")
	t.Execute(w, indexParams{
		app.TmplParams(c.App, r, "Index"),
	})
}
