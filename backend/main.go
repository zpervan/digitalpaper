package main

import (
	"digitalpaper/backend/core/logger"
	"digitalpaper/backend/web"
	"net/http"
)

func main() {
    logger.Info("Initializing dependencies")
    router := web.NewRoutes()

	logger.Info("Starting API server")
	err := http.ListenAndServe(":3500", router.HandleRequests())

	if err != nil {
		logger.Error("Error during server start-up. Reason: " + err.Error())
		return
	}
}
