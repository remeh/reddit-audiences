package main

import (
	"log"
	"os"

	"github.com/vrischmann/envconfig"
)

type App struct {
	db     Conn
	Config Config
}

type Config struct {
	DB string `envconfig:"DB,default=host=/var/run/postgresql sslmode=disable user=audiences dbname=audiences password=audiences"`
}

func (a *App) Init() {
	// init
	err := envconfig.Init(&a.Config)
	if err != nil {
		log.Println("err: on config reading:", err.Error())
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
