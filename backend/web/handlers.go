package web

import (
	"digitalpaper/backend/core/logger"

	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
)

// For local (non-Docker) development/testing
const localDatabaseUrl = "mongodb://admin:password@localhost:27018"

var database Database
var files []string
var fileServer http.Handler

func init() {
	// Initialize database
	var databaseUrl string

	if os.Getenv("TASKS_DB_ADDRESS") != "" {
		databaseUrl = os.Getenv("TASKS_DB_ADDRESS")
	} else {
		databaseUrl = localDatabaseUrl
	}

	_, err := Connect(databaseUrl)

	if err != nil {
		logger.Warn("Could not connect to database. Reason:" + err.Error())
		panic(err)
	}

	// Initialize web page and file server
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

func createPost() {
	// @TODO: Implement post creation
	logger.Info("Create post functionality not implemented")
}

func editPost() {
	// @TODO: Implement post editing/updating
	logger.Info("Edit/update post functionality not implemented")
}

func deletePost() {
	// @TODO: Implement post deletion
	logger.Info("Delete post functionality not implemented")
}

func getPosts() {
	// @TODO: Implement posts fetching
	logger.Info("Fetch posts functionality not implemented")
}

func createUser() {
	// @TODO: Implement user creation
	logger.Info("Create user functionality not implemented")
}

func editUser() {
	// @TODO: Implement user editing/updating
	logger.Info("Edit/update user functionality not implemented")
}

func deleteUser() {
	// @TODO: Implement user deletion
	logger.Info("Delete user functionality not implemented")
}

func getUsers() {
	// @TODO: Implement users fetching
	logger.Info("Fetch users functionality not implemented")
}

func HandleRequests() *mux.Router {
	router := mux.NewRouter()

	router.Path("/").HandlerFunc(home)
	router.PathPrefix("/").Handler(http.StripPrefix("/static", fileServer))

	return router
}
