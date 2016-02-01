// Reddit audiences crawler
// Rémy Mathieu © 2016
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/vrischmann/envconfig"
)

type App struct {
	db     Conn
	router *mux.Router
	Config Config
}

type Config struct {
	DB         string `envconfig:"DB,default=host=/var/run/postgresql sslmode=disable user=audiences dbname=audiences password=audiences"`
	PublicDir  string `envconfig:"DIR,default=static/"`
	ListenAddr string `envconfig:"ADDR,default=:9000"`
}

func (a *App) Init() {
	// init
	err := envconfig.Init(&a.Config)
	if err != nil {
		log.Println("err: on config reading:", err.Error())
		os.Exit(1)
	}

	// router
	a.router = mux.NewRouter()
	a.prepareStatic()
	http.Handle("/", a.router)

	// open pg connection
	a.db.Init(a.Config)
}

// DB opens a connection to PostgreSQL.
func (a *App) DB() Conn {
	return a.db
}

func (a *App) StartJobs() {
	StartCrawlingJob(a)
}

// Server
// ----------------------

func (a *App) Listen() error {
	// Starts listening.
	return http.ListenAndServe(a.Config.ListenAddr, nil)
}

func (a *App) AddApi(pattern string, handler http.Handler) {
	a.router.PathPrefix("/api").Subrouter().Handle(pattern, handler)
}

func (a *App) prepareStatic() {
	// Add the final route, the static assets and pages.
	a.router.PathPrefix("/").Handler(http.FileServer(http.Dir(a.Config.PublicDir)))
	log.Println("info: serving static from directory", a.Config.PublicDir)
}
