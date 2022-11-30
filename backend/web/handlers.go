package web

import (
	"digitalpaper/backend/core"
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
	logger.Info("Creating new post")
	var newPost Post
	err := json.NewDecoder(req.Body).Decode(&newPost)

	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while creating post"}
		errorResponse.Respond()
		return
	}

	context := req.Context()
	newPost.Id = uuid.New().String()
	err = h.Database.createPost(&context, &newPost)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while creating post"}
		errorResponse.Respond()
		return
	}

	core.ResponseSuccess(&w, "Post created successfully")
}

func (h *Handler) editPost(w http.ResponseWriter, req *http.Request) {
	logger.Info("Editing post")

	var updatedPost Post
	err := json.NewDecoder(req.Body).Decode(&updatedPost)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while updating post"}
		errorResponse.Respond()
		return
	}

	context := req.Context()
	err = h.Database.updatePost(context, &updatedPost)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while updating post"}
		errorResponse.Respond()
		return
	}

	core.ResponseSuccess(&w, "Post edited successfully")
}

func (h *Handler) deletePost(w http.ResponseWriter, req *http.Request) {
	postId := mux.Vars(req)["id"]
	logger.Info("Deleting post with Id " + postId)

	context := req.Context()
	err := h.Database.deletePost(context, postId)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while deleting post"}
		errorResponse.Respond()
		return
	}

	core.ResponseSuccess(&w, "Post deleted successfully")
}

func (h *Handler) getPosts(w http.ResponseWriter, req *http.Request) {
	logger.Info("Fetching all posts")

	context := req.Context()
	posts, err := h.Database.getAllPosts(&context)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting posts"}
		errorResponse.Respond()
		return
	}

	err = core.EncodeResponse(&w, &posts)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting posts"}
		errorResponse.Respond()
		return
	}
}

func (h *Handler) getPostById(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	logger.Info(fmt.Sprintf("Fetching post with ID %s", id))

	context := req.Context()
	post, err := h.Database.getPostById(&context, id)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while querying post"}
		errorResponse.Respond()
		return
	}

	err = core.EncodeResponse(&w, &post)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while querying post"}
		errorResponse.Respond()
		return
	}
}

func (h *Handler) createUser(w http.ResponseWriter, req *http.Request) {
	logger.Info("Creating new user")

	var newUser User
	err := json.NewDecoder(req.Body).Decode(&newUser)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while creating user"}
		errorResponse.Respond()
		return
	}

	err = h.Database.createUser(req.Context(), &newUser)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while creating user"}
		errorResponse.Respond()
		return
	}

	core.ResponseSuccess(&w, "User created successfully")
}

func (h *Handler) editUser(w http.ResponseWriter, req *http.Request) {
	logger.Info("Editing user")

	var updatedUser User
	err := json.NewDecoder(req.Body).Decode(&updatedUser)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while editing user"}
		errorResponse.Respond()
		return
	}

	context := req.Context()
	err = h.Database.updateUser(context, &updatedUser)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while editing user"}
		errorResponse.Respond()
		return
	}

	core.ResponseSuccess(&w, "User edited successfully")
}

func (h *Handler) deleteUser() {
	// @TODO: Implement user deletion
	logger.Info("Delete user functionality not implemented")
}

func (h *Handler) getUsers(w http.ResponseWriter, req *http.Request) {
	logger.Info("Fetching all users")

	context := req.Context()

	// @TODO: Add functionality to return a certain amount of user?
	users, err := h.Database.getUsers(&context, -1)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting users"}
		errorResponse.Respond()
		return
	}

	err = core.EncodeResponse(&w, &users)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting users"}
		errorResponse.Respond()
		return
	}
}

func (h *Handler) getUserByUsername(w http.ResponseWriter, req *http.Request) {
	username := mux.Vars(req)["username"]

	logger.Info("Getting user \"" + username + "\"")

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	user, err := h.Database.getUserByUsername(req.Context(), username)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting user"}
		errorResponse.Respond()
		return
	}

	err = core.EncodeResponse(&w, &user)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting user"}
		errorResponse.Respond()
		return
	}
}
