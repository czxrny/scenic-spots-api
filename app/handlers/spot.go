package handlers

import (
	"encoding/json"
	"net/http"
	"scenic-spots-api/app/database/repositories"
	"scenic-spots-api/app/logger"
	"scenic-spots-api/models"
	"strconv"
	"strings"
)

func getSpot(response http.ResponseWriter, request *http.Request) {
	queryParams := models.SpotQueryParams{
		Name:      request.URL.Query().Get("name"),
		Latitude:  request.URL.Query().Get("latitude"),
		Longitude: request.URL.Query().Get("longitude"),
		Radius:    request.URL.Query().Get("radius"),
		Category:  request.URL.Query().Get("category"),
	}

	found, err := repositories.FindSpotsByQueryParams(queryParams, request.Context())
	if err != nil {
		ErrorResponse(response, "500", err.Error(), http.StatusInternalServerError)
	}
	response.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(response).Encode(found); err != nil {
		ErrorResponse(response, "500", "Error while JSON encoding", http.StatusInternalServerError)
		return
	}
}

func addSpot(response http.ResponseWriter, request *http.Request) {
	logger.Info("Add a new spot.")
}

func getSpotById(response http.ResponseWriter, request *http.Request, id string) {
	logger.Info("Get spot by ID " + id)
}

func updateSpotById(response http.ResponseWriter, request *http.Request, id string) {
	logger.Info("Update spot by ID " + id)
}

func deleteSpotById(response http.ResponseWriter, request *http.Request, id string) {
	logger.Info("Delete spot by ID " + id)
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
	if _, err := strconv.Atoi(parts[2]); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	spotId := parts[2]

	if numberOfParts == 3 {
		switch method {
		case "GET":
			getSpotById(response, request, spotId)
		case "PUT":
			updateSpotById(response, request, spotId)
		case "DELETE":
			deleteSpotById(response, request, spotId)
		}
	} else if numberOfParts == 4 {
		spotElement := parts[3]
		switch spotElement {
		case "photo":
			Photo(response, request, spotId)
		case "review":
			Review(response, request, spotId)
		}
	} else {
		http.NotFound(response, request)
	}
}
