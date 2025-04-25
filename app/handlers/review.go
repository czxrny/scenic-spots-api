package handlers

import (
	"net/http"
	"scenic-spots-api/app/logger"
)

func getReview(response http.ResponseWriter, request *http.Request, id string) {
	logger.Info("Get reviews for the spot")
}

func addReview(response http.ResponseWriter, request *http.Request, id string) {
	logger.Info("Add review to the spot")
}

func deleteReview(response http.ResponseWriter, request *http.Request, spotId string, reviewId string) {
	logger.Info("Delete specified review")
}

func Review(response http.ResponseWriter, request *http.Request, id string) {
	method := request.Method

	switch method {
	case "GET":
		getReview(response, request, id)
	case "POST":
		addReview(response, request, id)
	case "DELETE":
		// check if there is reviewid specified! if not, delete all!
		deleteReview(response, request, id, "22")
	}
}
