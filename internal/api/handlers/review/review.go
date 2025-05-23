package review

import (
	"net/http"
	helpers "scenic-spots-api/internal/api/helpers"
	reviewRepo "scenic-spots-api/internal/database/repositories/review"
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

	params := models.ReviewQueryParams{
		SpotId: spotId,
		Limit:  request.URL.Query().Get("limit"),
	}

	found, err := reviewRepo.GetReviews(request.Context(), params)
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, found)
}

func addReview(response http.ResponseWriter, request *http.Request, spotId string) {
	var newReview models.NewReview
	if err := helpers.DecodeAndValidateRequestBody(request, &newReview); err != nil {
		helpers.ErrorResponse(response, "Error while decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	newReview.SpotId = spotId

	found, err := reviewRepo.AddReview(request.Context(), newReview)
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, found)
}

func deleteAllReviews(response http.ResponseWriter, request *http.Request, spotId string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	// can delete only if jwt states that the user is an admin.
	if err := helpers.CanEditAsset(request, ""); err != nil {
		helpers.ErrorResponse(response, "Authorization error: "+err.Error(), http.StatusUnauthorized)
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
		helpers.HandleRepoError(response, err)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, review)
}

func updateReviewById(response http.ResponseWriter, request *http.Request, id string) {
	found, err := reviewRepo.FindReviewById(request.Context(), id)
	if err != nil {
		helpers.HandleRepoError(response, err)
	}

	review := found[0]

	if err := helpers.CanEditAsset(request, review.AddedBy); err != nil {
		helpers.ErrorResponse(response, "Authorization error: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if err := helpers.DecodeAndValidateRequestBody(request, &review); err != nil {
		helpers.ErrorResponse(response, "Error while decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedReview, err := reviewRepo.UpdateReviewById(request.Context(), id, review)
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

	found, err := reviewRepo.FindReviewById(request.Context(), id)
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	review := found[0]

	if err := helpers.CanEditAsset(request, review.AddedBy); err != nil {
		helpers.ErrorResponse(response, "Authorization error: "+err.Error(), http.StatusUnauthorized)
		return
	}

	err = reviewRepo.DeleteReviewById(request.Context(), id)
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
