package web

import (
    "digitalpaper/backend/core"
    "github.com/gorilla/mux"
    "net/http"
)

type Routes struct {
    App *core.Application
    Handlers *Handler
}

func NewRoutes(app *core.Application) *Routes{
    routes := &Routes{}
    routes.App = app
    routes.Handlers = NewHandler(routes.App)

    return routes
}


func (r Routes)HandleRequests() *mux.Router {
	router := mux.NewRouter()

	// GET
    router.Path("/posts").Methods(http.MethodGet).HandlerFunc(r.Handlers.getPosts)
    router.Path("/posts/{id}").Methods(http.MethodGet).HandlerFunc(r.Handlers.getPostById)
    router.Path("/users").Methods(http.MethodGet).HandlerFunc(r.Handlers.getUsers)
    router.Path("/users/{username}").Methods(http.MethodGet).HandlerFunc(r.Handlers.getUserByUsername)

	// POST
    router.Path("/posts").Methods(http.MethodPost).HandlerFunc(r.Handlers.createPost)
    router.Path("/users").Methods(http.MethodPost).HandlerFunc(r.Handlers.createUser)

    // PUT
    router.Path("/posts").Methods(http.MethodPut).HandlerFunc(r.Handlers.editPost)
    router.Path("/users").Methods(http.MethodPut).HandlerFunc(r.Handlers.editUser)

    // DELETE
    router.Path("/posts/{id}").Methods(http.MethodDelete).HandlerFunc(r.Handlers.deletePost)
    router.Path("/users/{username}").Methods(http.MethodDelete).HandlerFunc(r.Handlers.deleteUser)

	return router
}
