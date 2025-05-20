package review

import (
	"encoding/json"
	"errors"
	"net/http"
	reviewRepo "scenic-spots-api/app/database/repositories/review"
	helpers "scenic-spots-api/internal/apihelpers"
	"scenic-spots-api/internal/repoerrors"
	"scenic-spots-api/models"
	"strings"

	"github.com/go-playground/validator/v10"
)

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
		case "DELETE":
			deleteAllReviews(response, request, spotId)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else if numberOfParts == 5 {
		reviewId := parts[4]
		if reviewId == "" {
			helpers.ErrorResponse(response, "Missing review ID", http.StatusBadRequest)
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
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else {
		response.WriteHeader(http.StatusNotFound)
	}
}

func getReview(response http.ResponseWriter, request *http.Request, spotId string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	params := models.ReviewQueryParams{
		SpotId: spotId,
		Limit:  request.URL.Query().Get("limit"),
	}

	found, err := reviewRepo.GetReviews(request.Context(), params)
	if err != nil {
		if errors.Is(err, repoerrors.ErrDoesNotExist) {
			helpers.ErrorResponse(response, "Spot with ID: ["+spotId+"] was not found", http.StatusBadRequest)
			return
		}
		if errors.Is(err, repoerrors.ErrInvalidQueryParameters) {
			helpers.ErrorResponse(response, err.Error(), http.StatusBadRequest)
			return
		}
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(response).Encode(found); err != nil {
		helpers.ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func addReview(response http.ResponseWriter, request *http.Request, spotId string) {
	var newReview models.NewReview
	if err := json.NewDecoder(request.Body).Decode(&newReview); err != nil {
		helpers.ErrorResponse(response, "Bad request body", http.StatusBadRequest)
		return
	}
	newReview.SpotId = spotId

	validate := validator.New()
	if err := validate.Struct(newReview); err != nil {
		helpers.ErrorResponse(response, "Invalid parameters", http.StatusBadRequest)
		return
	}

	found, err := reviewRepo.AddReview(request.Context(), newReview)
	if err != nil {
		if errors.Is(err, repoerrors.ErrDoesNotExist) {
			helpers.ErrorResponse(response, "Spot with ID: ["+spotId+"] was not found", http.StatusBadRequest)
			return
		}
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(response).Encode(found); err != nil {
		helpers.ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteAllReviews(response http.ResponseWriter, request *http.Request, spotId string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	if err := reviewRepo.DeleteAllReviews(request.Context(), spotId); err != nil {
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func getReviewById(response http.ResponseWriter, request *http.Request, id string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	review, err := reviewRepo.FindReviewById(request.Context(), id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrDoesNotExist) {
			helpers.ErrorResponse(response, "Review with ID: ["+id+"] was not found", http.StatusBadRequest)
			return
		}
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(response).Encode(review); err != nil {
		helpers.ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateReviewById(response http.ResponseWriter, request *http.Request, id string) {
	var review models.NewReview
	if err := json.NewDecoder(request.Body).Decode(&review); err != nil {
		helpers.ErrorResponse(response, "Bad request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(review); err != nil {
		helpers.ErrorResponse(response, "Invalid parameters", http.StatusBadRequest)
		return
	}

	updatedSpot, err := reviewRepo.UpdateReviewById(request.Context(), id, review)
	if err != nil {
		if errors.Is(err, repoerrors.ErrDoesNotExist) {
			helpers.ErrorResponse(response, "Review with ID: ["+id+"] was not found", http.StatusBadRequest)
			return
		}
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(response).Encode(updatedSpot); err != nil {
		helpers.ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteReviewById(response http.ResponseWriter, request *http.Request, id string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	err := reviewRepo.DeleteReviewById(request.Context(), id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrDoesNotExist) {
			helpers.ErrorResponse(response, "Review with ID: ["+id+"] was not found", http.StatusNotFound)
			return
		}
		helpers.ErrorResponse(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
