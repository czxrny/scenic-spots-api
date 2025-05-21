package user

import (
	"net/http"
	helpers "scenic-spots-api/internal/api/helpers"
	"scenic-spots-api/internal/auth"
	userRepo "scenic-spots-api/internal/database/repositories/user"
	"scenic-spots-api/internal/models"
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
		helpers.HandleRepoError(response, err)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, result)
}

func loginUser(response http.ResponseWriter, request *http.Request) {
	var userCredentials models.UserCredentials
	if err := helpers.DecodeAndValidateRequestBody(request, &userCredentials); err != nil {
		helpers.ErrorResponse(response, "Error while decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := auth.LoginUser(request.Context(), userCredentials)
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, result)
}

func deleteUserById(response http.ResponseWriter, request *http.Request, userId string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "DELETE request must not contain a body", http.StatusBadRequest)
		return
	}

	err := userRepo.DeleteUserById(request.Context(), userId)
	if err != nil {
		helpers.HandleRepoError(response, err)
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
