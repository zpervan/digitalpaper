package web

import (
	"digitalpaper/backend/core"
	"digitalpaper/backend/database"
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
	App      *core.Application
	Database *database.Database
}

func NewHandler(app *core.Application) *Handler {
	var databaseUrl string

	if os.Getenv("TASKS_DB_ADDRESS") != "" {
		databaseUrl = os.Getenv("TASKS_DB_ADDRESS")
	} else {
		databaseUrl = localDatabaseUrl
	}

	databaseTemp, err := database.NewDatabase(app, databaseUrl)

	if err != nil {
		app.Log.Warn("Could not connect to database. Reason: " + err.Error())
		panic(err)
	}

	handler := &Handler{}
	handler.App = app
	handler.Database = databaseTemp

	return handler
}

func (h *Handler) createPost(w http.ResponseWriter, req *http.Request) {
	h.App.Log.Info("Creating new post")
	var newPost core.Post
	err := json.NewDecoder(req.Body).Decode(&newPost)

	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while creating post"}
		errorResponse.Respond()
		return
	}

	context := req.Context()
	newPost.Id = uuid.New().String()
	err = h.Database.CreatePost(&context, &newPost)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while creating post"}
		errorResponse.Respond()
		return
	}

	core.ResponseSuccess(&w, "Post created successfully")
}

func (h *Handler) editPost(w http.ResponseWriter, req *http.Request) {
	h.App.Log.Info("Editing post")

	var updatedPost core.Post
	err := json.NewDecoder(req.Body).Decode(&updatedPost)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while updating post"}
		errorResponse.Respond()
		return
	}

	context := req.Context()
	err = h.Database.UpdatePost(context, &updatedPost)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while updating post"}
		errorResponse.Respond()
		return
	}

	core.ResponseSuccess(&w, "Post edited successfully")
}

func (h *Handler) deletePost(w http.ResponseWriter, req *http.Request) {
	postId := mux.Vars(req)["id"]
	h.App.Log.Info("Deleting post with Id " + postId)

	context := req.Context()
	err := h.Database.DeletePost(context, postId)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while deleting post"}
		errorResponse.Respond()
		return
	}

	core.ResponseSuccess(&w, "Post deleted successfully")
}

func (h *Handler) getPosts(w http.ResponseWriter, req *http.Request) {
	h.App.Log.Info("Fetching all posts")

	context := req.Context()
	posts, err := h.Database.GetAllPosts(&context)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting posts"}
		errorResponse.Respond()
		return
	}

	err = core.EncodeResponse(&w, http.StatusOK, &posts)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting posts"}
		errorResponse.Respond()
		return
	}
}

func (h *Handler) getPostById(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	h.App.Log.Info(fmt.Sprintf("Fetching post with ID %s", id))

	context := req.Context()
	post, err := h.Database.GetPostById(&context, id)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while querying post"}
		errorResponse.Respond()
		return
	}

	err = core.EncodeResponse(&w, http.StatusOK, &post)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while querying post"}
		errorResponse.Respond()
		return
	}
}

func (h *Handler) createUser(w http.ResponseWriter, req *http.Request) {
	h.App.Log.Info("Creating new user")

	var newUser core.User
	err := json.NewDecoder(req.Body).Decode(&newUser)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while creating user"}
		errorResponse.Respond()
		return
	}

	validationResult := core.ValidateUser(&newUser)
	if validationResult != nil {
		err = core.EncodeResponse(&w, http.StatusNotAcceptable, validationResult)
		if err != nil {
			h.App.Log.Error("cannot encode validation results. reason: " + err.Error())
			return
		}

		validationError := "error while validating input data during user creation:"

		for _, value := range validationResult {
			validationError += "\n" + value.Attribute + " - " + value.Error
		}

		h.App.Log.Warn(validationError)

		return
	}

	userExists, err := h.Database.UserExists(req.Context(), &newUser)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while checking user existence"}
		errorResponse.Respond()

		h.App.Log.Error("error while checking user existence. reason: " + err.Error())
		return
	}

	if userExists {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("user with that username or mail already exits"))
		return
	}

	err = h.Database.CreateUser(req.Context(), &newUser)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while creating user"}
		errorResponse.Respond()
		return
	}

	core.ResponseSuccess(&w, "User created successfully")
}

func (h *Handler) editUser(w http.ResponseWriter, req *http.Request) {
	h.App.Log.Info("Editing user")

	var updatedUser core.User
	err := json.NewDecoder(req.Body).Decode(&updatedUser)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while editing user"}
		errorResponse.Respond()
		return
	}

	context := req.Context()
	err = h.Database.UpdateUser(context, &updatedUser)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while editing user"}
		errorResponse.Respond()
		return
	}

	core.ResponseSuccess(&w, "User edited successfully")
}

func (h *Handler) deleteUser(w http.ResponseWriter, req *http.Request) {
	username := mux.Vars(req)["username"]
	h.App.Log.Info("Deleting user \"" + username + "\"")

	context := req.Context()
	err := h.Database.DeleteUser(context, username)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while deleting user"}
		errorResponse.Respond()
		return
	}

	core.ResponseSuccess(&w, "User deleted successfully")
}

func (h *Handler) getUsers(w http.ResponseWriter, req *http.Request) {
	h.App.Log.Info("Fetching all users")

	context := req.Context()

	// @TODO: Add functionality to return a certain amount of user?
	users, err := h.Database.GetUsers(&context, -1)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting users"}
		errorResponse.Respond()
		return
	}

	err = core.EncodeResponse(&w, http.StatusOK, &users)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting users"}
		errorResponse.Respond()
		return
	}
}

func (h *Handler) getUserByUsername(w http.ResponseWriter, req *http.Request) {
	username := mux.Vars(req)["username"]

	h.App.Log.Info("Getting user \"" + username + "\"")

	user, err := h.Database.GetUserByUsername(req.Context(), username)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting user"}
		errorResponse.Respond()
		return
	}

	err = core.EncodeResponse(&w, http.StatusOK, &user)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting user"}
		errorResponse.Respond()
		return
	}
}
