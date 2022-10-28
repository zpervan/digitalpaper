package main

// @TODO: Add custom logger class - severity, date, time and log file

import (
	"digitalpaper/backend/web"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello, Digital Paper!")

	err := http.ListenAndServe("localhost:3500", web.HandleRequests())

	if err != nil {
		fmt.Println("Error while creating server. Reason: ", err.Error())
		return
	}
}
