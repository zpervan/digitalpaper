package main

import (
	"digitalpaper/backend/core/logger"
	"digitalpaper/backend/web"
	"net/http"
)

func main() {
	logger.Info("Starting API server...")
	err := http.ListenAndServe(":3500", web.HandleRequests())

	if err != nil {
		logger.Error("Error during server start-up. Reason: " + err.Error())
		return
	}
}
