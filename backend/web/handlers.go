package web

import (
	"digitalpaper/backend/core/logger"
	"encoding/json"
    "fmt"
    "github.com/gorilla/mux"

    "net/http"
	"os"

	"github.com/google/uuid"
)

// For local (non-Docker) development/testing
const localDatabaseUrl = "mongodb://admin:password@localhost:27018"

type Handler struct {
	Database *Database
}

func NewHandler() *Handler {
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

    handler := &Handler{}
	handler.Database = &databaseTemp

	return handler
}

func (h *Handler) createPost(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/posts" {
		http.NotFound(w, req)
		return
	}

	// @TODO: Find a more elegant way to enable the CORS policy
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var newPost Post
	err := json.NewDecoder(req.Body).Decode(&newPost)

	if err != nil {
		logger.Error("Could not create new post. Reason:" + err.Error())
	}

	newPost.Id = uuid.New().String()

	err = h.Database.createPost(&newPost)

	if err != nil {
		logger.Error("Could not create a new task in the database. Reason:" + err.Error())
	}
}

func (h *Handler) editPost() {
	// @TODO: Implement post editing/updating
	logger.Info("Edit/update post functionality not implemented")
}

func (h *Handler) deletePost() {
	// @TODO: Implement post deletion
	logger.Info("Delete post functionality not implemented")
}

func (h *Handler) getPosts(w http.ResponseWriter, req *http.Request) {
	logger.Info("Fetching all posts")

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	context := req.Context()

	posts, err := h.Database.getAllPosts(&context)
	if err != nil {
		logger.Error(err.Error())
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *Handler)getPostById(w http.ResponseWriter, req *http.Request) {
    id := mux.Vars(req)["id"]
    logger.Info(fmt.Sprintf("Fetching post with ID %s", id))

    context := req.Context()

    post, err := h.Database.getPostById(&context, id)
    if err != nil {
        errorMessage := fmt.Sprintf("%d - error while querying post by ID. %s", http.StatusInternalServerError, err.Error())

        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(errorMessage))

        return
    }

    err = json.NewEncoder(w).Encode(post)
    if err != nil {
        errorMessage := fmt.Sprintf("%d - error while querying post by ID. %s", http.StatusInternalServerError, err.Error())
        logger.Error(errorMessage)

        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(errorMessage))

        return
    }
}

func (h *Handler) createUser(w http.ResponseWriter, req *http.Request) {
    if req.URL.Path != "/users" {
        http.NotFound(w, req)
        return
    }

    logger.Info("Creating new user")

    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    var newUser User

    err := json.NewDecoder(req.Body).Decode(&newUser)
    if err != nil {
        errorMessage := fmt.Sprintf("%d - error while creating new user. %s", http.StatusInternalServerError, err.Error())
        logger.Error(errorMessage)

        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(errorMessage))

        return
    }

    err = h.Database.createUser(req.Context(), &newUser)
    if err != nil {
        errorMessage := fmt.Sprintf("%d - error while creating new user. %s", http.StatusInternalServerError, err.Error())
        logger.Error(errorMessage)

        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(errorMessage))

        return
    }

}

func (h *Handler) editUser() {
	// @TODO: Implement user editing/updating
	logger.Info("Edit/update user functionality not implemented")
}

func (h *Handler) deleteUser() {
	// @TODO: Implement user deletion
	logger.Info("Delete user functionality not implemented")
}

func (h *Handler) getUsers() {
	// @TODO: Implement users fetching
	logger.Info("Fetch users functionality not implemented")
}
