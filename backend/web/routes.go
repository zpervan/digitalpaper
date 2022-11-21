package web

import (
	"github.com/gorilla/mux"
)

func HandleRequests() *mux.Router {
	router := mux.NewRouter()

	// GET
	router.Path("/posts").Methods("GET").HandlerFunc(getPosts)

	// POST
	router.Path("/posts").Methods("POST").HandlerFunc(createPost)

	return router
}
