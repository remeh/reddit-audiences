// Reddit audiences crawler
// Rémy Mathieu © 2016
package main

import (
	. "github.com/remeh/reddit-audiences/api"
	"github.com/remeh/reddit-audiences/app"
	"github.com/remeh/reddit-audiences/web"
)

func main() {
	var a app.App
	a.Init()
	a.StartJobs()

	declareWebRoutes(&a)
	declareApiRoutes(&a)
	a.Listen()
}

func declareWebRoutes(a *app.App) {
	// Finally index
	a.Add("/audiences/{subreddit}", web.Audiences{a}, "GET")

	a.Add("/signup", web.SignupGet{a}, "GET")
	a.Add("/signup", web.SignupPost{a}, "POST")

	a.Add("/index", web.Index{a}, "GET")
	a.Add("/", web.Index{a}, "GET")
}

func declareApiRoutes(a *app.App) {
	a.AddApi("/today/{subreddit}", LogRoute(a, TodayHandler{a}), "GET")
}
