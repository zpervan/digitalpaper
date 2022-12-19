package web

import (
	"digitalpaper/backend/core"
	"fmt"
	"net/http"
)

type Middleware struct {
	App *core.Application
}

func NewMiddleware(app *core.Application) *Middleware {
	middleware := &Middleware{}

	middleware.App = app

	return middleware
}

func (m *Middleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !m.App.SessionManager.Exists(req.Context(), "authenticatedUserId") {
			errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: fmt.Errorf(""), StatusCode: http.StatusForbidden}
			errorResponse.Respond()

			m.App.Log.Error(errorResponse.Message)

			return
		}

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, req)
	})
}
