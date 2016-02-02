// Reddit audiences crawler
// Rémy Mathieu © 2016
package app

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/vrischmann/envconfig"
)

type App struct {
	db        Conn
	router    *mux.Router
	Templates *template.Template
	Config    Config
}

type Config struct {
	DB           string `envconfig:"DB,default=host=/var/run/postgresql sslmode=disable user=audiences dbname=audiences password=audiences"`
	PublicDir    string `envconfig:"DIR,default=static/"`
	Crawl        bool   `envconfig:"CRAWL,default=true`
	TemplatesDir string `envconfig:"TEMPLATES,default=templates/"`
	ListenAddr   string `envconfig:"ADDR,default=:9000"`
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

	// read templates
	if err := a.initTemplates(); err != nil {
		log.Println("err: can't read templates:", err.Error())
		os.Exit(1)
	}

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
	a.prepareStatic()

	http.Handle("/", a.router)

	// Starts listening.
	return http.ListenAndServe(a.Config.ListenAddr, nil)
}

func (a *App) Add(pattern string, handler http.Handler) {
	a.router.Handle(pattern, handler)
}

func (a *App) AddApi(pattern string, handler http.Handler) {
	a.router.PathPrefix("/api").Subrouter().Handle(pattern, handler)
}

func (a *App) prepareStatic() {
	// Add the final route, the static assets and pages.
	a.router.PathPrefix("/").Handler(http.FileServer(http.Dir(a.Config.PublicDir)))
	log.Println("info: serving static from directory", a.Config.PublicDir)
}

func (a *App) initTemplates() error {
	templates, err := ReadTemplates(a)

	if err != nil {
		return err
	}

	a.Templates = templates
	log.Println("info: using templates from the directory", a.Config.TemplatesDir)
	return nil
}
