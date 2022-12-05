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
}

func NewRoutes(app *core.Application) *Routes {
	routes := &Routes{}
	routes.App = app
	routes.Handlers = NewHandler(routes.App)

	return routes
}

func (r Routes) HandleRequests() *mux.Router {
	router := mux.NewRouter()

    dynamicMiddleware := alice.New(r.App.SessionManager.LoadAndSave)
    
	// GET
	router.Path("/posts").Methods(http.MethodGet).Handler(dynamicMiddleware.ThenFunc(r.Handlers.getPosts))
	router.Path("/posts/{id}").Methods(http.MethodGet).Handler(dynamicMiddleware.ThenFunc(r.Handlers.getPostById))
	router.Path("/users").Methods(http.MethodGet).Handler(dynamicMiddleware.ThenFunc(r.Handlers.getUsers))
	router.Path("/users/{username}").Methods(http.MethodGet).Handler(dynamicMiddleware.ThenFunc(r.Handlers.getUserByUsername))

	// POST
	router.Path("/posts").Methods(http.MethodPost).Handler(dynamicMiddleware.ThenFunc(r.Handlers.createPost))
	router.Path("/users").Methods(http.MethodPost).Handler(dynamicMiddleware.ThenFunc(r.Handlers.createUser))

	// PUT
	router.Path("/posts").Methods(http.MethodPut).Handler(dynamicMiddleware.ThenFunc(r.Handlers.editPost))
	router.Path("/users").Methods(http.MethodPut).Handler(dynamicMiddleware.ThenFunc(r.Handlers.editUser))

	// DELETE
	router.Path("/posts/{id}").Methods(http.MethodDelete).Handler(dynamicMiddleware.ThenFunc(r.Handlers.deletePost))
	router.Path("/users/{username}").Methods(http.MethodDelete).Handler(dynamicMiddleware.ThenFunc(r.Handlers.deleteUser))

    // @TODO: Add middleware functionalities (i.e. security headers)
	return router
}
