package tests

import (
	"digitalpaper/backend/core"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ErrorResponseTestSuite struct {
	er core.ErrorResponse
}

func Test_GivenValidErrorResponseDataWithRaisedError_WhenResponding_ThenErrorResponseIsCorrect(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	er := core.ErrorResponse{
		ResponseWriter: responseRecorder,
		RaisedError:    fmt.Errorf("some silly error"),
		StatusCode:     http.StatusBadRequest,
		Message:        "This is an error message",
	}

	er.Respond()

	expectedStatusCode := http.StatusBadRequest
	actualStatusCode := responseRecorder.Result().StatusCode

	assert.Equal(t, expectedStatusCode, actualStatusCode)

	expectedErrorMessage := "400 - This is an error message. reason: some silly error"
	// @TODO: HeaderMap is deprecated, find a more suitable solution
	actualErrorMessage := responseRecorder.HeaderMap.Get("X-Status-Reason")

	assert.Equal(t, expectedErrorMessage, actualErrorMessage)
}

func Test_GivenValidErrorResponseDataWithoutRaisedError_WhenResponding_ThenErrorResponseIsCorrect(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	er := core.ErrorResponse{
		ResponseWriter: responseRecorder,
		RaisedError:    nil,
		StatusCode:     http.StatusForbidden,
		Message:        "This error is wow",
	}

	er.Respond()

	expectedStatusCode := http.StatusForbidden
	actualStatusCode := responseRecorder.Result().StatusCode

	assert.Equal(t, expectedStatusCode, actualStatusCode)

	expectedErrorMessage := "403 - This error is wow."
	actualErrorMessage := responseRecorder.HeaderMap.Get("X-Status-Reason")

	assert.Equal(t, expectedErrorMessage, actualErrorMessage)
}

func Test_GivenInvalidErrorResponse_WhenResponding_ThenPanicIsTriggered(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	er := core.ErrorResponse{
		ResponseWriter: responseRecorder,
		RaisedError:    nil,
		StatusCode:     http.StatusForbidden,
		Message:        "",
	}

	assert.Panics(t, er.Respond)
}

func Test_GivenValidData_WhenResponding_ThenResponseIsCorrect(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	err := core.EncodeResponse(responseRecorder, http.StatusOK, "some dummy data")
	assert.Nil(t, err)

	expectedStatusCode := http.StatusOK
	actualStatusCode := responseRecorder.Result().StatusCode
	assert.Equal(t, expectedStatusCode, actualStatusCode)

	expectedData := "\"some dummy data\"\n"
	actualData := responseRecorder.Body.String()
	assert.Equal(t, expectedData, actualData)
}

func Test_GivenInvalidData_WhenResponding_ThenErrorIsRaised(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	// We pass an unsupported/invalid data type to the JSON serializer which will raise an error
	err := core.EncodeResponse(responseRecorder, http.StatusOK, func() {})
	assert.NotNil(t, err)
}
