package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/emiljannesson/bed-and-breakfast/internal/config"
	"github.com/emiljannesson/bed-and-breakfast/internal/handlers"
	"github.com/emiljannesson/bed-and-breakfast/internal/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager

func main() {
	// change to true when in production
	appConfig.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache, err:", err)
	}
	appConfig.TemplateCache = tc

	appConfig.UseCache = false
	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)

	render.NewTemplates(&appConfig)

	fmt.Printf("Starting application on port %s", portNumber)

	srv := &http.Server{Addr: portNumber, Handler: routes(&appConfig)}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
