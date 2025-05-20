package user

import (
	"net/http"
	"scenic-spots-api/app/logger"
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
	logger.Info("User register attempt.")
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
