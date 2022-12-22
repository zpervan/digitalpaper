package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/stretchr/testify/suite"

	"digitalpaper/backend/core"
	"digitalpaper/backend/core/logger"
	"digitalpaper/backend/web"
)

type MiddlewareTestSuite struct {
	suite.Suite
	middleware *web.Middleware
	serverMock *TestServer
}

// Suite functions

func (mts *MiddlewareTestSuite) SetupSuite() {
	fmt.Println("Setting up middleware test suite...")

	mockApp := &core.Application{Log: logger.New(), SessionManager: scs.New()}
	mockApp.SessionManager.Lifetime = 12 * time.Hour

	mts.middleware = web.NewMiddleware(mockApp)

	// Dummy router
	mockRouter := mux.NewRouter()

	// Middleware mock
	mw := alice.New(mts.middleware.App.SessionManager.LoadAndSave)
	protectedMw := mw.Append(mts.middleware.RequireAuthentication)

	// Dummy paths
	mockRouter.Path(dummyGetOkUrl).Methods(http.MethodGet).Handler(protectedMw.ThenFunc(DummyGetOkResponse))
	mockRouter.Path(dummyLogin).Methods(http.MethodGet).Handler(mw.ThenFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx = context.WithValue(ctx, "authenticatedUserId", "0000")
		req = req.WithContext(ctx)

		mts.middleware.App.SessionManager.Put(req.Context(), "authenticatedUserId", "0000")
	}))
	mockRouter.Path(removeUser).Methods(http.MethodGet).Handler(mw.ThenFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx = context.WithValue(ctx, "authenticatedUserId", "0000")
		req = req.WithContext(ctx)

		// @TODO: Handle error
		mts.middleware.App.SessionManager.Destroy(req.Context())
	}))

	mts.serverMock = NewTestServer(mts.T(), mockRouter)

	fmt.Println("Setting up middleware test suite... COMPLETE")
}

func (mts *MiddlewareTestSuite) TearDownSuite() {
	fmt.Println("Tearing down middleware test suite...")

	mts.serverMock.Close()

	fmt.Println("Tearing down middleware test suite... COMPLETE")
}

func (mts *MiddlewareTestSuite) TearDownTest() {
	mts.serverMock.ExecuteGet(mts.T(), removeUser)
}

// Tests

func (mts *MiddlewareTestSuite) Test_GivenValidHttpRequest_WhenReturningResponse_ThenSecureHeaderDataIsCorrect() {
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mts.middleware.SecureHeaders(next).ServeHTTP(responseRecorder, request)

	result := responseRecorder.Result()

	expected := "origin-when-cross-origin"
	actual := result.Header.Get("Referrer-Policy")
	mts.Assert().Equal(expected, actual)

	expected = "nosniff"
	actual = result.Header.Get("X-Content-Type-Options")
	mts.Assert().Equal(expected, actual)

	expected = "deny"
	actual = result.Header.Get("X-Frame-Options")
	mts.Assert().Equal(expected, actual)

	expected = "0"
	actual = result.Header.Get("X-XSS-Protection")
	mts.Assert().Equal(expected, actual)
}

func (mts *MiddlewareTestSuite) Test_GivenARequestFromNonAuthorizedUser_WhenRequestIsProcessed_ThenRequestIsDeclined() {
	result, _ := mts.serverMock.ExecuteGet(mts.T(), dummyGetOkUrl)

	expected := http.StatusForbidden
	actual := result.StatusCode
	mts.Assertions.Equal(expected, actual, "status code should br \"403 - Forbidden\", but isn't")
}

func (mts *MiddlewareTestSuite) Test_GivenARequestFromAuthorizedUser_WhenRequestIsProcessed_ThenRequestPasses() {
	resultLogin, _ := mts.serverMock.ExecuteGet(mts.T(), dummyLogin)

	expected := http.StatusOK
	actual := resultLogin.StatusCode
	mts.Assertions.Equal(expected, actual, "status code should br \"200 - OK\", but isn't")

	resultGet, _ := mts.serverMock.ExecuteGet(mts.T(), dummyGetOkUrl)

	expected = http.StatusOK
	actual = resultGet.StatusCode
	mts.Assertions.Equal(expected, actual, "status code should br \"200 - OK\", but isn't")
}

func TestRunMiddleware(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}
