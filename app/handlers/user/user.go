package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"scenic-spots-api/app/auth"
	userRepo "scenic-spots-api/app/database/repositories/user"
	helpers "scenic-spots-api/internal/apihelpers"
	"scenic-spots-api/internal/repoerrors"
	"scenic-spots-api/models"
	"scenic-spots-api/utils/logger"
	"strings"
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
		default:
			UserById(response, request, operation)
		}
	} else {
		response.WriteHeader(http.StatusNotFound)
	}
}

func UserById(response http.ResponseWriter, request *http.Request, id string) {
	method := request.Method

	switch method {
	case "DELETE":
		if err := helpers.IsAuthenticated(request); err != nil {
			helpers.ErrorResponse(response, err.Error(), http.StatusUnauthorized)
			return
		}
		if err := helpers.CanEditAsset(request, id); err != nil {
			helpers.ErrorResponse(response, err.Error(), http.StatusUnauthorized)
			return
		}
		deleteUserById(response, request, id)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func registerUser(response http.ResponseWriter, request *http.Request) {
	var userRegisterInfo models.UserRegisterInfo
	if err := helpers.DecodeAndValidateRequestBody(request, &userRegisterInfo); err != nil {
		helpers.ErrorResponse(response, "Error while decoding request body: "+err.Error(), http.StatusBadRequest)
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
	if err := helpers.DecodeAndValidateRequestBody(request, &userCredentials); err != nil {
		helpers.ErrorResponse(response, "Error while decoding request body: "+err.Error(), http.StatusBadRequest)
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

func deleteUserById(response http.ResponseWriter, request *http.Request, userId string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "DELETE request must not contain a body", http.StatusBadRequest)
		return
	}

	err := userRepo.DeleteUserById(request.Context(), userId)
	if err != nil {
		if errors.Is(err, repoerrors.ErrDoesNotExist) {
			helpers.ErrorResponse(response, "User does not exist in database.", http.StatusNotFound)
			return
		}
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

/* TO IMPLEMENT IN THE FUTURE - ADDITIONAL USER INFO, STORED IN DIFFERENT COLLECTION */

func getUserInfo(response http.ResponseWriter, request *http.Request) {
	logger.Info("Fetching user profile info.")
}

func setUserPassword(response http.ResponseWriter, request *http.Request) {
	logger.Info("Changing user password.")
}

func updateUserInfo(response http.ResponseWriter, request *http.Request) {
	logger.Info("Updating user profile info.")
}
