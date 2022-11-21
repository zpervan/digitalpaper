package web

import (
	"digitalpaper/backend/core/logger"
	"encoding/json"

	"net/http"
	"os"

	"github.com/google/uuid"
)

// For local (non-Docker) development/testing
const localDatabaseUrl = "mongodb://admin:password@localhost:27018"

var database *Database

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
}

func createPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/posts" {
		http.NotFound(w, r)
		return
	}

    // @TODO: Find a more elegant way to enable the CORS policy
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
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

    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

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
