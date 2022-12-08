package main

import (
	"digitalpaper/backend/core"
	"digitalpaper/backend/core/logger"
	"digitalpaper/backend/web"
	"net/http"
	"time"

	"github.com/alexedwards/scs/mongodbstore"
	"github.com/alexedwards/scs/v2"
	"github.com/rs/cors"
)

func main() {
	logger := logger.New()
	logger.Info("Initializing dependencies")

	// Core application functionalities initialization
	app := &core.Application{}
	app.Log = logger

	app.SessionManager = scs.New()
	app.SessionManager.Lifetime = 12 * time.Hour

	router := web.NewRoutes(app)
	app.SessionManager.Store = mongodbstore.New(router.Handlers.Database.Sessions.Database())

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	})

	// Starting server functionalities
	app.Log.Info("Starting API server")
	err := http.ListenAndServe(":3500", cors.Handler(router.HandleRequests()))

	if err != nil {
		app.Log.Error("Error during server start-up. Reason: " + err.Error())
		return
	}
}
