package web

import (
	"digitalpaper/backend/core"
	"digitalpaper/backend/database"
	"encoding/json"
	"fmt"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

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

	if os.Getenv("DB_URL") != "" {
		databaseUrl = os.Getenv("DB_URL")
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

		h.App.Log.Error(errorResponse.Message)
		return
	}

	context := req.Context()
	newPost.Id = uuid.New().String()
	err = h.Database.CreatePost(&context, &newPost)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while creating post"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
		return
	}
}

func (h *Handler) editPost(w http.ResponseWriter, req *http.Request) {
	h.App.Log.Info("Editing post")

	var updatedPost core.Post
	err := json.NewDecoder(req.Body).Decode(&updatedPost)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while updating post"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
		return
	}

	context := req.Context()
	err = h.Database.UpdatePost(context, &updatedPost)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while updating post"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
		return
	}
}

func (h *Handler) deletePost(w http.ResponseWriter, req *http.Request) {
	postId := mux.Vars(req)["id"]
	h.App.Log.Info("Deleting post with Id " + postId)

	context := req.Context()
	err := h.Database.DeletePost(context, postId)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while deleting post"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
		return
	}
}

func (h *Handler) getPosts(w http.ResponseWriter, req *http.Request) {
	h.App.Log.Info("Fetching all posts")

	context := req.Context()
	posts, err := h.Database.GetAllPosts(&context)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting posts"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
		return
	}

	err = core.EncodeResponse(&w, http.StatusOK, &posts)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting posts"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
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

		h.App.Log.Error(errorResponse.Message)
		return
	}

	err = core.EncodeResponse(&w, http.StatusOK, &post)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while querying post"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
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

		h.App.Log.Error(errorResponse.Message)
		return
	}

	validationResult := core.ValidateUser(&newUser)
	if validationResult != nil {
		err = core.EncodeResponse(&w, http.StatusBadRequest, validationResult)
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

		h.App.Log.Error(errorResponse.Message)
		return
	}

	if userExists {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusBadRequest, Message: "user with that username or mail already exists"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
		return
	}

	err = h.Database.CreateUser(req.Context(), &newUser)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while creating user"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
		return
	}
}

func (h *Handler) editUser(w http.ResponseWriter, req *http.Request) {
	h.App.Log.Info("Editing user")

	var updatedUser core.User
	err := json.NewDecoder(req.Body).Decode(&updatedUser)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while editing user"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
		return
	}

	context := req.Context()
	err = h.Database.UpdateUser(context, &updatedUser)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while editing user"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
		return
	}
}

func (h *Handler) deleteUser(w http.ResponseWriter, req *http.Request) {
	username := mux.Vars(req)["username"]
	h.App.Log.Info("Deleting user \"" + username + "\"")

	context := req.Context()
	err := h.Database.DeleteUser(context, username)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while deleting user"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
		return
	}
}

func (h *Handler) getUsers(w http.ResponseWriter, req *http.Request) {
	h.App.Log.Info("Fetching all users")

	context := req.Context()

	// @TODO: Add functionality to return a certain amount of user?
	users, err := h.Database.GetUsers(&context, -1)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting users"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
		return
	}

	err = core.EncodeResponse(&w, http.StatusOK, &users)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting users"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
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

		h.App.Log.Error(errorResponse.Message)
		return
	}

	err = core.EncodeResponse(&w, http.StatusOK, &user)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while getting user"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
		return
	}
}

func (h *Handler) login(w http.ResponseWriter, req *http.Request) {
	h.App.Log.Info("Logging in")

	var userCredentials core.UserLogin
	err := json.NewDecoder(req.Body).Decode(&userCredentials)
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusInternalServerError, Message: "error while login"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
		return
	}

	fetchedUser, err := h.Database.GetUserByMail(req.Context(), userCredentials.Mail)
	if err != nil {
		statusCode := http.StatusInternalServerError

		// If the data is empty, the user doesn't exist
		if fetchedUser.IsEmpty() {
			statusCode = http.StatusNotFound
		}

		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: nil, StatusCode: statusCode, Message: "error while login"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(userCredentials.Password))
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: fmt.Errorf("wrong credentials"), StatusCode: http.StatusNotAcceptable, Message: ""}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.RaisedError.Error())
		return
	}

	err = h.App.SessionManager.RenewToken(req.Context())
	if err != nil {
		errorResponse := core.ErrorResponse{ResponseWriter: &w, RaisedError: err, StatusCode: http.StatusNotAcceptable, Message: "could not create session for user"}
		errorResponse.Respond()

		h.App.Log.Error(errorResponse.Message)
		return
	}

	h.App.SessionManager.Put(req.Context(), "user", fetchedUser.Username)
	h.App.Log.Info("User logged in")
}
