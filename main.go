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

	a.Add("/register", web.RegisterGet{a}, "GET")
	a.Add("/register", web.RegisterPost{a}, "POST")

	a.Add("/signin", web.SigninGet{a}, "GET")
	a.Add("/signin", web.SigninPost{a}, "POST")

	a.Add("/account", web.Account{a}, "GET")
	a.Add("/logout", web.Logout{a}, "GET")

	a.Add("/index", web.Index{a}, "GET")
	a.Add("/", web.Index{a}, "GET")
}

func declareApiRoutes(a *app.App) {
	a.AddApi("/today/{subreddit}", LogRoute(a, TodayHandler{a}), "GET")
}
