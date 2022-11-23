package web

import (
    "github.com/gorilla/mux"
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
	router.Path("/posts").Methods("GET").HandlerFunc(r.Handlers.getPosts)
    router.Path("/posts/{id}").Methods("GET").HandlerFunc(r.Handlers.getPostById)

	// POST
    router.Path("/posts").Methods("POST").HandlerFunc(r.Handlers.createPost)

	return router
}
