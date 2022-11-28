package web

import (
    "github.com/gorilla/mux"
    "net/http"
)

type Routes struct {
    Handlers *Handler
}

func NewRoutes() *Routes{
    routes := &Routes{}
    routes.Handlers = NewHandler()

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

	return router
}
