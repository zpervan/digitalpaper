package main

// @TODO: Add custom logger class - severity, date, time and log file

import (
	"digitalpaper/backend/core/logger"
	"digitalpaper/backend/web"
	"net/http"
)

func main() {
	logger.Info("Starting API server...")
	err := http.ListenAndServe("localhost:3500", web.HandleRequests())

	if err != nil {
		logger.Error("Error while creating server. Reason: " + err.Error())
		return
	}
}
