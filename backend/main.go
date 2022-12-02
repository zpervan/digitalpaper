package main

import (
	"digitalpaper/backend/core"
	"digitalpaper/backend/core/logger"
	"digitalpaper/backend/web"
	"net/http"
)

func main() {
	app := &core.Application{}
	app.Log = logger.New()

    app.Log.Info("Initializing dependencies")
    router := web.NewRoutes(app)

	app.Log.Info("Starting API server")
	err := http.ListenAndServe(":3500", router.HandleRequests())

	if err != nil {
		app.Log.Error("Error during server start-up. Reason: " + err.Error())
		return
	}
}
