package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// @TODO: Recursive and automatic HTML file assignment
	files := []string{"./ui/html/base.html"}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func HandleRequests() *mux.Router {
	router := mux.NewRouter()

	router.Path("/").HandlerFunc(home)

	return router
}
