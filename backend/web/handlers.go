package web

import (
	"digitalpaper/backend/core/logger"
	"encoding/json"

	"html/template"
	"net/http"
	"os"

	"github.com/google/uuid"
)

// For local (non-Docker) development/testing
const localDatabaseUrl = "mongodb://admin:password@localhost:27018"

var database *Database
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

	databaseTemp, err := Connect(databaseUrl)

	if err != nil {
		logger.Warn("Could not connect to database. Reason: " + err.Error())
		panic(err)
	}

	database = &databaseTemp

	// Initialize web page and file server
	// @TODO: Populate automatically HTML file list
	files = []string{
		"./ui/html/pages/base.html",
		"./ui/html/components/header_with_image.html",
		"./ui/html/components/navigation_bar.html",
		"./ui/html/components/preview_article.html",
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

	err = ts.ExecuteTemplate(w, "base", &postMock)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/posts" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var newPost Post
	err := json.NewDecoder(r.Body).Decode(&newPost)

	if err != nil {
		logger.Error("Could not create new post. Reason:" + err.Error())
	}

	newPost.Id = uuid.New().String()

	err = database.CreatePost(&newPost)

	if err != nil {
		logger.Error("Could not create a new task in the database. Reason:" + err.Error())
	}
}

func editPost() {
	// @TODO: Implement post editing/updating
	logger.Info("Edit/update post functionality not implemented")
}

func deletePost() {
	// @TODO: Implement post deletion
	logger.Info("Delete post functionality not implemented")
}

func getPosts(w http.ResponseWriter, req *http.Request) {
	logger.Info("Fetching all posts...")

	context := req.Context()

	posts, err := database.GetAllPosts(&context)
	if err != nil {
		logger.Error(err.Error())
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		logger.Error(err.Error())
	}
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
