package web

import (
	"digitalpaper/backend/core/logger"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var files []string
var fileServer http.Handler

func init() {
	// @TODO: Populate automatically HTML file list
	files = []string{
		"./ui/html/pages/base.html",
		"./ui/html/components/navigation_bar.html",
	}

	fileServer = http.FileServer(http.Dir("./ui/static/"))
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func HandleRequests() *mux.Router {
	router := mux.NewRouter()

	router.Path("/").HandlerFunc(home)
	router.PathPrefix("/").Handler(http.StripPrefix("/static", fileServer))

	return router
}
