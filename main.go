// Reddit audiences crawler
// Rémy Mathieu © 2016
package main

func main() {
	var app App
	app.Init()
	go app.StartJobs()

	declareApiRoutes(&app)
	app.Listen()
}

func declareApiRoutes(a *App) {
	a.AddApi("/today/{subreddit}", LogRoute(a, TodayHandler{a}))
}
