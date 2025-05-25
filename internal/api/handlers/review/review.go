package review

import (
	"net/http"
	helpers "scenic-spots-api/internal/api/helpers"
	reviewService "scenic-spots-api/internal/api/service/review"
	"scenic-spots-api/internal/models"
	"strings"
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
			if err := helpers.IsAuthenticated(request); err != nil {
				helpers.ErrorResponse(response, err.Error(), http.StatusUnauthorized)
				return
			}
			addReview(response, request, spotId)
		case "DELETE":
			if err := helpers.IsAuthenticated(request); err != nil {
				helpers.ErrorResponse(response, err.Error(), http.StatusUnauthorized)
				return
			}
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
			if err := helpers.IsAuthenticated(request); err != nil {
				helpers.ErrorResponse(response, err.Error(), http.StatusUnauthorized)
				return
			}
			updateReviewById(response, request, reviewId)
		case "DELETE":
			if err := helpers.IsAuthenticated(request); err != nil {
				helpers.ErrorResponse(response, err.Error(), http.StatusUnauthorized)
				return
			}
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

	found, err := reviewService.GetReview(request.Context(), spotId, request.URL.Query())
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, found)
}

func addReview(response http.ResponseWriter, request *http.Request, spotId string) {
	token, err := helpers.GetJWTToken(request)
	if err != nil {
		helpers.ErrorResponse(response, "Error while decoding header: "+err.Error(), http.StatusBadRequest)
		return
	}

	var newReview models.NewReview
	if err := helpers.DecodeAndValidateRequestBody(request, &newReview); err != nil {
		helpers.ErrorResponse(response, "Error while decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	newReview.SpotId = spotId

	found, err := reviewService.AddReview(request.Context(), token, newReview)
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, found)
}

func getReviewById(response http.ResponseWriter, request *http.Request, id string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	review, err := reviewService.FindReviewById(request.Context(), id)
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, review)
}

func updateReviewById(response http.ResponseWriter, request *http.Request, id string) {
	token, err := helpers.GetJWTToken(request)
	if err != nil {
		helpers.ErrorResponse(response, "Error while decoding header: "+err.Error(), http.StatusBadRequest)
		return
	}

	var reviewInfo models.NewReview
	if err := helpers.DecodeAndValidateRequestBody(request, &reviewInfo); err != nil {
		helpers.ErrorResponse(response, "Error while decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedReview, err := reviewService.UpdateReviewById(request.Context(), token, reviewInfo, id)
	if err != nil {
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, updatedReview)
}

func deleteReviewById(response http.ResponseWriter, request *http.Request, id string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}
	token, err := helpers.GetJWTToken(request)
	if err != nil {
		helpers.ErrorResponse(response, "Error while decoding header: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = reviewService.DeleteReviewById(request.Context(), token, id)
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func deleteAllReviews(response http.ResponseWriter, request *http.Request, spotId string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}
	token, err := helpers.GetJWTToken(request)
	if err != nil {
		helpers.ErrorResponse(response, "Error while decoding header: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := reviewService.DeleteAllReviews(request.Context(), token, spotId); err != nil {
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
