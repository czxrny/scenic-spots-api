package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"scenic-spots-api/app/database/repositories"
	"scenic-spots-api/models"
	"strings"

	"github.com/go-playground/validator/v10"
)

func getSpot(response http.ResponseWriter, request *http.Request) {
	queryParams := models.SpotQueryParams{
		Name:      request.URL.Query().Get("name"),
		Latitude:  request.URL.Query().Get("latitude"),
		Longitude: request.URL.Query().Get("longitude"),
		Radius:    request.URL.Query().Get("radius"),
		Category:  request.URL.Query().Get("category"),
	}

	found, err := repositories.GetSpot(queryParams, request.Context())
	if err != nil {
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

func addSpot(response http.ResponseWriter, request *http.Request) {
	var spot models.NewSpot
	if err := json.NewDecoder(request.Body).Decode(&spot); err != nil {
		ErrorResponse(response, "Bad request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(spot); err != nil {
		ErrorResponse(response, "Invalid parameters", http.StatusBadRequest)
		return
	}

	ctx := request.Context()
	addedSpot, err := repositories.AddSpot(spot, ctx)
	if err != nil {
		if errors.Is(err, repositories.ErrSpotAlreadyExists) {
			ErrorResponse(response, "Spot already exists in the database", http.StatusConflict)
			return
		}
		ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(response).Encode(addedSpot); err != nil {
		ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getSpotById(response http.ResponseWriter, request *http.Request, id string) {
	ctx := request.Context()

	spot, err := repositories.FindSpotById(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrDoesNotExist) {
			ErrorResponse(response, "Spot with ID: ["+id+"] was not found", http.StatusNotFound)
			return
		}
		ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(response).Encode(spot); err != nil {
		ErrorResponse(response, "Failed to encode JSON error response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateSpotById(response http.ResponseWriter, request *http.Request, id string) {
	var spot models.NewSpot
	if err := json.NewDecoder(request.Body).Decode(&spot); err != nil {
		ErrorResponse(response, "Bad request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(spot); err != nil {
		ErrorResponse(response, "Invalid parameters", http.StatusBadRequest)
		return
	}

	ctx := request.Context()
	updatedSpot, err := repositories.UpdateSpot(ctx, id, spot)
	if err != nil {
		if errors.Is(err, repositories.ErrDoesNotExist) {
			ErrorResponse(response, "Spot with ID: ["+id+"] was not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, repositories.ErrSpotAlreadyExists) {
			ErrorResponse(response, "Spot in this coordinates already exists!", http.StatusConflict)
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

func deleteSpotById(response http.ResponseWriter, request *http.Request, id string) {
	ctx := request.Context()

	err := repositories.DeleteSpotById(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrDoesNotExist) {
			ErrorResponse(response, "Spot with ID: ["+id+"] was not found", http.StatusNotFound)
			return
		}
		ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func Spot(response http.ResponseWriter, request *http.Request) {
	method := request.Method
	switch method {
	case "GET":
		getSpot(response, request)
	case "POST":
		addSpot(response, request)
	}
}

func SpotById(response http.ResponseWriter, request *http.Request) {
	parts := strings.Split(request.URL.Path, "/")
	numberOfParts := len(parts)
	method := request.Method
	spotId := parts[2]

	if numberOfParts == 3 {
		switch method {
		case "GET":
			getSpotById(response, request, spotId)
		case "PATCH":
			updateSpotById(response, request, spotId)
		case "DELETE":
			deleteSpotById(response, request, spotId)
		default:
			response.WriteHeader(http.StatusNotFound)
		}
	} else if numberOfParts >= 4 {
		spotElement := parts[3]
		switch spotElement {
		case "photo":
			Photo(response, request, spotId)
		case "review":
			Review(response, request, spotId)
		default:
			response.WriteHeader(http.StatusNotFound)
		}
	}
}
