package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"scenic-spots-api/app/database/repositories"
	"scenic-spots-api/models"
	"strings"

	"github.com/go-playground/validator/v10"
)

func getReview(response http.ResponseWriter, request *http.Request, spotId string) {
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		ErrorResponse(response, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	if len(bodyBytes) > 0 {
		ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	ctx := request.Context()
	params := models.ReviewQueryParams{
		SpotId: spotId,
		Limit:  request.URL.Query().Get("limit"),
	}

	found, err := repositories.GetReviews(ctx, params)
	if err != nil {
		if errors.Is(err, repositories.ErrDoesNotExist) {
			ErrorResponse(response, "Spot with ID: ["+spotId+"] was not found", http.StatusBadRequest)
			return
		}
		if errors.Is(err, repositories.ErrInvalidQueryParameters) {
			ErrorResponse(response, err.Error(), http.StatusBadRequest)
			return
		}
		ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(response).Encode(found); err != nil {
		ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func addReview(response http.ResponseWriter, request *http.Request, spotId string) {
	ctx := request.Context()

	var newReview models.NewReview
	if err := json.NewDecoder(request.Body).Decode(&newReview); err != nil {
		ErrorResponse(response, "Bad request body", http.StatusBadRequest)
		return
	}
	newReview.SpotId = spotId

	validate := validator.New()
	if err := validate.Struct(newReview); err != nil {
		ErrorResponse(response, "Invalid parameters", http.StatusBadRequest)
		return
	}

	found, err := repositories.AddReview(ctx, newReview)

	if err != nil {
		if errors.Is(err, repositories.ErrDoesNotExist) {
			ErrorResponse(response, "Spot with ID: ["+spotId+"] was not found", http.StatusBadRequest)
			return
		}
		ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(response).Encode(found); err != nil {
		ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getReviewById(response http.ResponseWriter, request *http.Request, id string) {
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		ErrorResponse(response, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	if len(bodyBytes) > 0 {
		ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	ctx := request.Context()
	review, err := repositories.FindReviewById(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrDoesNotExist) {
			ErrorResponse(response, "Review with ID: ["+id+"] was not found", http.StatusBadRequest)
			return
		}
		ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(response).Encode(review); err != nil {
		ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateReviewById(response http.ResponseWriter, request *http.Request, id string) {
	var review models.NewReview
	if err := json.NewDecoder(request.Body).Decode(&review); err != nil {
		ErrorResponse(response, "Bad request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(review); err != nil {
		ErrorResponse(response, "Invalid parameters", http.StatusBadRequest)
		return
	}

	ctx := request.Context()
	updatedSpot, err := repositories.UpdateReviewById(ctx, id, review)
	if err != nil {
		if errors.Is(err, repositories.ErrDoesNotExist) {
			ErrorResponse(response, "Review with ID: ["+id+"] was not found", http.StatusBadRequest)
			return
		}
		ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(response).Encode(updatedSpot); err != nil {
		ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteReviewById(response http.ResponseWriter, request *http.Request, id string) {
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		ErrorResponse(response, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	if len(bodyBytes) > 0 {
		ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	ctx := request.Context()
	err = repositories.DeleteReviewById(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrDoesNotExist) {
			ErrorResponse(response, "Review with ID: ["+id+"] was not found", http.StatusNotFound)
			return
		}
		ErrorResponse(response, err.Error(), http.StatusInternalServerError)
		return
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
		if reviewId == "" {
			ErrorResponse(response, "Missing review ID", http.StatusBadRequest)
			return
		}
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
	} else {
		response.WriteHeader(http.StatusNotFound)
	}
}
