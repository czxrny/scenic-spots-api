package handlers

import (
	"encoding/json"
	"net/http"
	"scenic-spots-api/app/database/repositories"
	"scenic-spots-api/app/logger"
	"strings"
)

func getReview(response http.ResponseWriter, request *http.Request, spotId string) {
	ctx := request.Context()
	limit := request.URL.Query().Get("limit")

	found, err := repositories.GetReviews(ctx, spotId, limit)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			response.WriteHeader(http.StatusNotFound)
		}
		ErrorResponse(response, "500", err.Error(), http.StatusInternalServerError)
	}

	response.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(response).Encode(found); err != nil {
		ErrorResponse(response, "500", "Error while JSON encoding", http.StatusInternalServerError)
		return
	}
}

func addReview(response http.ResponseWriter, request *http.Request, spotId string) {
	logger.Info("Add review to the spot")
}

func deleteReview(response http.ResponseWriter, request *http.Request, spotId string, reviewId string) {
	logger.Info("Delete specified review")
}

func Review(response http.ResponseWriter, request *http.Request, spotId string) {
	method := request.Method

	switch method {
	case "GET":
		getReview(response, request, spotId)
	case "POST":
		addReview(response, request, spotId)
	case "DELETE":
		// check if there is reviewid specified! if not, delete all!
		deleteReview(response, request, spotId, "22")
	}
}
