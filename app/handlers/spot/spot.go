package spot

import (
	"encoding/json"
	"errors"
	"net/http"
	spotRepo "scenic-spots-api/app/database/repositories/spot"
	pHandler "scenic-spots-api/app/handlers/photo"
	rHandler "scenic-spots-api/app/handlers/review"
	helpers "scenic-spots-api/internal/apihelpers"
	"scenic-spots-api/internal/repoerrors"
	"scenic-spots-api/models"
	"strings"
)

func Spot(response http.ResponseWriter, request *http.Request) {
	method := request.Method
	switch method {
	case "GET":
		getSpot(response, request)
	case "POST":
		if err := helpers.IsAuthenticated(request); err != nil {
			helpers.ErrorResponse(response, err.Error(), http.StatusUnauthorized)
			return
		}
		addSpot(response, request)
	}
}

func SpotById(response http.ResponseWriter, request *http.Request) {
	parts := strings.Split(request.URL.Path, "/")
	numberOfParts := len(parts)
	method := request.Method
	spotId := parts[2]

	if spotId == "" {
		helpers.ErrorResponse(response, "Missing spot ID", http.StatusBadRequest)
		return
	}

	if numberOfParts == 3 {
		switch method {
		case "GET":
			getSpotById(response, request, spotId)
		case "PATCH":
			if err := helpers.IsAuthenticated(request); err != nil {
				helpers.ErrorResponse(response, err.Error(), http.StatusUnauthorized)
				return
			}
			updateSpotById(response, request, spotId)
		case "DELETE":
			if err := helpers.IsAuthenticated(request); err != nil {
				helpers.ErrorResponse(response, err.Error(), http.StatusUnauthorized)
				return
			}
			deleteSpotById(response, request, spotId)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else if numberOfParts >= 4 {
		spotElement := parts[3]
		switch spotElement {
		case "photo":
			pHandler.Photo(response, request, spotId)
		case "review":
			rHandler.Review(response, request, spotId)
		default:
			response.WriteHeader(http.StatusNotFound)
		}
	} else {
		response.WriteHeader(http.StatusNotFound)
	}
}

func getSpot(response http.ResponseWriter, request *http.Request) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	queryParams := models.SpotQueryParams{
		Name:      request.URL.Query().Get("name"),
		Latitude:  request.URL.Query().Get("latitude"),
		Longitude: request.URL.Query().Get("longitude"),
		Radius:    request.URL.Query().Get("radius"),
		Category:  request.URL.Query().Get("category"),
	}

	found, err := spotRepo.GetSpot(request.Context(), queryParams)
	if err != nil {
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

func addSpot(response http.ResponseWriter, request *http.Request) {
	var spot models.NewSpot
	if err := helpers.DecodeAndValidateRequestBody(request, &spot); err != nil {
		helpers.ErrorResponse(response, "Error while decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	addedSpot, err := spotRepo.AddSpot(request.Context(), spot)
	if err != nil {
		if errors.Is(err, repoerrors.ErrAlreadyExists) {
			helpers.ErrorResponse(response, "Spot already exists in the database", http.StatusConflict)
			return
		}
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(response).Encode(addedSpot); err != nil {
		helpers.ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getSpotById(response http.ResponseWriter, request *http.Request, id string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	spot, err := spotRepo.FindSpotById(request.Context(), id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrDoesNotExist) {
			helpers.ErrorResponse(response, "Spot with ID: ["+id+"] was not found", http.StatusNotFound)
			return
		}
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(response).Encode(spot); err != nil {
		helpers.ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateSpotById(response http.ResponseWriter, request *http.Request, id string) {
	var spot models.NewSpot
	if err := helpers.DecodeAndValidateRequestBody(request, &spot); err != nil {
		helpers.ErrorResponse(response, "Error while decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedSpot, err := spotRepo.UpdateSpot(request.Context(), id, spot)
	if err != nil {
		if errors.Is(err, repoerrors.ErrDoesNotExist) {
			helpers.ErrorResponse(response, "Spot with ID: ["+id+"] was not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, repoerrors.ErrAlreadyExists) {
			helpers.ErrorResponse(response, "Spot in this coordinates already exists!", http.StatusConflict)
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

func deleteSpotById(response http.ResponseWriter, request *http.Request, id string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	err := spotRepo.DeleteSpotById(request.Context(), id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrDoesNotExist) {
			helpers.ErrorResponse(response, "Spot with ID: ["+id+"] was not found", http.StatusNotFound)
			return
		}
		helpers.ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
