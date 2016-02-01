// Reddit audiences crawler
// Rémy Mathieu © 2016
package main

func main() {
	var app App
	app.Init()
	go app.StartJobs()
	app.Listen()
}

func declareApiRoutes(a *App) {

}
