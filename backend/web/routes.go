package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

func HandleRequests() *mux.Router {
	router := mux.NewRouter()

	// GET
	router.Path("/posts").Methods("GET").HandlerFunc(getPosts)

	// POST
	router.Path("/posts").Methods("POST").HandlerFunc(createPost)

	// Web page
	router.Path("/").HandlerFunc(home)
	router.PathPrefix("/").Handler(http.StripPrefix("/static", fileServer))

	return router
}
