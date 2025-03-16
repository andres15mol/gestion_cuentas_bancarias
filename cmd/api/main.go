package main

import (
	"fmt"
	"gestion_cuentas_bancarias/internal/data"
	"log"
	"net/http"
	"os"
)

type config struct {
	port int
}

type application struct {
	config     config
	infoLog    *log.Logger
	errorLog   *log.Logger
	models     data.Models
	enviroment string
}

func main() {
	var cfg config
	cfg.port = 8080

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	environment := os.Getenv("ENV")

	app := &application{
		config:     cfg,
		infoLog:    infoLog,
		errorLog:   errorLog,
		models:     data.New(),
		enviroment: environment,
	}

	err := app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) serve() error {
	app.infoLog.Println("API listening on port", app.config.port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	return srv.ListenAndServe()
}
