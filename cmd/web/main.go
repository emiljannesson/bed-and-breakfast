package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/emiljannesson/bed-and-breakfast/internal/config"
	"github.com/emiljannesson/bed-and-breakfast/internal/driver"
	"github.com/emiljannesson/bed-and-breakfast/internal/handlers"
	"github.com/emiljannesson/bed-and-breakfast/internal/helpers"
	"github.com/emiljannesson/bed-and-breakfast/internal/models"
	"github.com/emiljannesson/bed-and-breakfast/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println("Starting application on port", portNumber)

	srv := &http.Server{Addr: portNumber, Handler: routes(&appConfig)}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// change to true when in production
	appConfig.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=") // temporarily hard coded
	if err != nil {
		log.Fatal("Cannot connect to database, exiting...")
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		return db, err
	}
	appConfig.TemplateCache = tc
	appConfig.UseCache = false

	repo := handlers.NewRepo(&appConfig, db)
	handlers.NewHandlers(repo)
	render.NewTemplates(&appConfig)
	helpers.NewHelpers(&appConfig)

	return db, nil
}
