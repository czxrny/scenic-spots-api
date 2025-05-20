package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"scenic-spots-api/app/auth"
	"scenic-spots-api/app/logger"
	helpers "scenic-spots-api/internal/apihelpers"
	"scenic-spots-api/internal/repoerrors"
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
	var userRegisterInfo models.UserRegisterInfo
	if err := json.NewDecoder(request.Body).Decode(&userRegisterInfo); err != nil {
		helpers.ErrorResponse(response, "Bad request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(userRegisterInfo); err != nil {
		helpers.ErrorResponse(response, "Invalid parameters", http.StatusBadRequest)
		return
	}

	result, err := auth.RegisterUser(request.Context(), userRegisterInfo)
	if err != nil {
		if errors.Is(err, repoerrors.ErrAlreadyExists) {
			helpers.ErrorResponse(response, "User already exists in the database", http.StatusConflict)
			return
		}
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(response).Encode(result); err != nil {
		helpers.ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func loginUser(response http.ResponseWriter, request *http.Request) {
	var userCredentials models.UserCredentials
	if err := json.NewDecoder(request.Body).Decode(&userCredentials); err != nil {
		helpers.ErrorResponse(response, "Bad request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(userCredentials); err != nil {
		helpers.ErrorResponse(response, "Invalid parameters", http.StatusBadRequest)
		return
	}

	result, err := auth.LoginUser(request.Context(), userCredentials)
	if err != nil {
		if errors.Is(err, repoerrors.ErrDoesNotExist) {
			helpers.ErrorResponse(response, "User does not exist in database.", http.StatusConflict)
			return
		}
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(response).Encode(result); err != nil {
		helpers.ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
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
