package handlers

import (
	"encoding/json"
	"net/http"
	"scenic-spots-api/app/database/repositories"
	"scenic-spots-api/models"
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
	ctx := request.Context()

	var newReview models.NewReview
	if err := json.NewDecoder(request.Body).Decode(&newReview); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	found, err := repositories.AddReview(ctx, spotId, newReview)

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

func getReviewById(response http.ResponseWriter, request *http.Request, id string) {
	ctx := request.Context()

	spot, err := repositories.FindReviewById(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			response.WriteHeader(http.StatusNotFound)
			return
		}
		ErrorResponse(response, "500", err.Error(), http.StatusInternalServerError)
	}

	response.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(response).Encode(spot); err != nil {
		ErrorResponse(response, "500", "Error while JSON encoding", http.StatusInternalServerError)
		return
	}
}

func updateReviewById(response http.ResponseWriter, request *http.Request, id string) {
	var review models.NewReview
	if err := json.NewDecoder(request.Body).Decode(&review); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := request.Context()
	updatedSpot, err := repositories.UpdateReviewById(ctx, id, review)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			response.WriteHeader(http.StatusConflict)
			return
		}
		ErrorResponse(response, "500", "Error while adding the spot to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(response).Encode(updatedSpot); err != nil {
		ErrorResponse(response, "500", "Error while JSON encoding", http.StatusInternalServerError)
		return
	}
}

func deleteReviewById(response http.ResponseWriter, request *http.Request, id string) {
	ctx := request.Context()

	err := repositories.DeleteReviewById(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			response.WriteHeader(http.StatusNotFound)
			return
		}
		ErrorResponse(response, "500", err.Error(), http.StatusInternalServerError)
	}

	response.WriteHeader(http.StatusNoContent)
}

func Review(response http.ResponseWriter, request *http.Request, spotId string) {
	parts := strings.Split(request.URL.Path, "/")
	numberOfParts := len(parts)
	method := request.Method

	if numberOfParts == 4 {
		switch method {
		case "GET":
			getReview(response, request, spotId)
		case "POST":
			addReview(response, request, spotId)
		default:
			response.WriteHeader(http.StatusNotFound)
		}
	} else if numberOfParts == 5 {
		reviewId := parts[4]
		switch method {
		case "GET":
			getReviewById(response, request, reviewId)
		case "PATCH":
			updateReviewById(response, request, reviewId)
		case "DELETE":
			deleteReviewById(response, request, reviewId)
		default:
			response.WriteHeader(http.StatusNotFound)
		}
	}
}
