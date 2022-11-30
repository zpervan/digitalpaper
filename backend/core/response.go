package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	"digitalpaper/backend/core/logger"
)

type ErrorResponse struct {
	ResponseWriter   *http.ResponseWriter
	RaisedError error
	StatusCode  int
	Message     string
}

func (er *ErrorResponse) Respond() {
	errorMessage := fmt.Sprintf("%d - %s. %s", er.StatusCode, er.Message, er.RaisedError.Error())
	logger.Error(errorMessage)

	(*er.ResponseWriter).WriteHeader(http.StatusInternalServerError)
	(*er.ResponseWriter).Write([]byte(errorMessage))
}

func ResponseSuccess(w *http.ResponseWriter, message string) {
	(*w).WriteHeader(http.StatusOK)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set(CorsString, AllowedAddressForCors)
	(*w).Write([]byte(message))
}

func EncodeResponse(w *http.ResponseWriter, data any) error {
	(*w).Header().Set(CorsString, AllowedAddressForCors)
	err := json.NewEncoder(*w).Encode(data)

	if err != nil {
		return err
	}

	return nil
}
