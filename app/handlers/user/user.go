package user

import (
	"encoding/json"
	"net/http"
	"scenic-spots-api/app/auth"
	userRepo "scenic-spots-api/app/database/repositories/user"
	"scenic-spots-api/app/logger"
	helpers "scenic-spots-api/internal/apihelpers"
	"scenic-spots-api/models"
	"strings"

	"github.com/go-playground/validator/v10"
)

func User(response http.ResponseWriter, request *http.Request) {
	parts := strings.Split(request.URL.Path, "/")
	numberOfParts := len(parts)
	method := request.Method

	if numberOfParts <= 3 {
		operation := parts[2]
		switch operation {
		case "register":
			if method != "POST" {
				response.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			registerUser(response, request)

		case "login":
			if method != "POST" {
				response.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			loginUser(response, request)

		case "me":
			if method != "GET" {
				response.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			getUserInfo(response, request)

		case "password":
			if method != "PUT" {
				response.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			setUserPassword(response, request)

		case "profile":
			if method != "GET" {
				response.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			updateUserInfo(response, request)

		case "delete":
			if method != "DELETE" {
				response.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			deleteUser(response, request)
		default:
			response.WriteHeader(http.StatusNotFound)
		}
	} else {
		response.WriteHeader(http.StatusNotFound)
	}

}

func registerUser(response http.ResponseWriter, request *http.Request) {
	var newUser models.UserCredentials
	if err := json.NewDecoder(request.Body).Decode(&newUser); err != nil {
		helpers.ErrorResponse(response, "Bad request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(newUser); err != nil {
		helpers.ErrorResponse(response, "Invalid parameters", http.StatusBadRequest)
		return
	}

	addedUser, err := userRepo.AddUser(request.Context(), newUser)
	if err != nil {
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := auth.CreateToken(addedUser[0])
	if err != nil {
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result := models.UserTokenResponse{
		Token:   token,
		LocalId: addedUser[0].Id,
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(response).Encode(result); err != nil {
		helpers.ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func loginUser(response http.ResponseWriter, request *http.Request) {
	logger.Info("User login attempt.")
}

func getUserInfo(response http.ResponseWriter, request *http.Request) {
	logger.Info("Fetching user profile info.")
}

func setUserPassword(response http.ResponseWriter, request *http.Request) {
	logger.Info("Changing user password.")
}

func updateUserInfo(response http.ResponseWriter, request *http.Request) {
	logger.Info("Updating user profile info.")
}

func deleteUser(response http.ResponseWriter, request *http.Request) {
	logger.Info("Deleting user account.")
}
