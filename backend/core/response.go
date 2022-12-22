package core

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	ResponseWriter http.ResponseWriter
	RaisedError    error
	StatusCode     int
	Message        string
}

func (er *ErrorResponse) Respond() {
	if (er.Message == "") && (er.RaisedError == nil) {
		panic("provide a descriptive error message to the error response")
	}

	var errorMessage string

	if er.RaisedError == nil {
		errorMessage = fmt.Sprintf("%d - %s.", er.StatusCode, er.Message)
	} else {
		errorMessage = fmt.Sprintf("%d - %s. reason: %s", er.StatusCode, er.Message, er.RaisedError.Error())
	}

	er.ResponseWriter.WriteHeader(er.StatusCode)
	er.ResponseWriter.Header().Set("X-Status-Reason", errorMessage)
}

func EncodeResponse(w http.ResponseWriter, statusCode int, data any) error {
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		return err
	}

	return nil
}
