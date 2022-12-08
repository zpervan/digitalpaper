package core

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	ResponseWriter   *http.ResponseWriter
	RaisedError error
	StatusCode  int
	Message     string
}

func (er *ErrorResponse) Respond() {
	if er.RaisedError == nil {
		er.RaisedError = fmt.Errorf("")
	}

	errorMessage := fmt.Sprintf("%d - %s. %s", er.StatusCode, er.Message, er.RaisedError.Error())

	(*er.ResponseWriter).WriteHeader(er.StatusCode)
	(*er.ResponseWriter).Header().Set("X-Status-Reason", errorMessage)
}

func EncodeResponse(w *http.ResponseWriter, statusCode int, data any) error {
	(*w).WriteHeader(statusCode)

	err := json.NewEncoder(*w).Encode(data)

	if err != nil {
		return err
	}

	return nil
}
