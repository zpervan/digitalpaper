package web

import (
	"digitalpaper/backend/core"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

type Routes struct {
	App      *core.Application
	Handlers *Handler
	Middleware *Middleware
}

func NewRoutes(app *core.Application) *Routes {
	routes := &Routes{}
	routes.App = app
	routes.Handlers = NewHandler(routes.App)
	routes.Middleware = NewMiddleware(routes.App)

	return routes
}

func (r Routes) HandleRequests() http.Handler {
	router := mux.NewRouter()
	
	// Unprotected
    dynamicMw := alice.New(r.App.SessionManager.LoadAndSave)
    
	// GET
	router.Path("/posts").Methods(http.MethodGet).Handler(dynamicMw.ThenFunc(r.Handlers.getPosts))
	router.Path("/posts/{id}").Methods(http.MethodGet).Handler(dynamicMw.ThenFunc(r.Handlers.getPostById))

	// POST
	router.Path("/users").Methods(http.MethodPost).Handler(dynamicMw.ThenFunc(r.Handlers.createUser))
	router.Path("/login").Methods(http.MethodPost).Handler(dynamicMw.ThenFunc(r.Handlers.login))

	// Protected
	protectedMw := dynamicMw.Append(r.Middleware.RequireAuthentication)

	// GET
	router.Path("/users").Methods(http.MethodGet).Handler(protectedMw.ThenFunc(r.Handlers.getUsers))
	router.Path("/users/{username}").Methods(http.MethodGet).Handler(protectedMw.ThenFunc(r.Handlers.getUserByUsername))

	// POST
	router.Path("/posts").Methods(http.MethodPost).Handler(protectedMw.ThenFunc(r.Handlers.createPost))
	router.Path("/logout").Methods(http.MethodPost).Handler(protectedMw.ThenFunc(r.Handlers.logout))

	// PUT
	router.Path("/posts").Methods(http.MethodPut).Handler(protectedMw.ThenFunc(r.Handlers.editPost))
	router.Path("/users").Methods(http.MethodPut).Handler(protectedMw.ThenFunc(r.Handlers.editUser))

	// DELETE
	router.Path("/posts/{id}").Methods(http.MethodDelete).Handler(protectedMw.ThenFunc(r.Handlers.deletePost))
	router.Path("/users/{username}").Methods(http.MethodDelete).Handler(protectedMw.ThenFunc(r.Handlers.deleteUser))

	// This will be called before each handler
	standardMwChain := alice.New(r.Middleware.SecureHeaders)

	return standardMwChain.Then(router)
}
