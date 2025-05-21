package spot

import (
	"net/http"
	pHandler "scenic-spots-api/internal/api/handlers/photo"
	rHandler "scenic-spots-api/internal/api/handlers/review"
	helpers "scenic-spots-api/internal/api/helpers"
	spotRepo "scenic-spots-api/internal/database/repositories/spot"
	"scenic-spots-api/internal/models"
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
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
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
		helpers.HandleRepoError(response, err)
		return
	}
	helpers.WriteJSONResponse(response, http.StatusOK, found)
}

func addSpot(response http.ResponseWriter, request *http.Request) {
	var spot models.NewSpot
	if err := helpers.DecodeAndValidateRequestBody(request, &spot); err != nil {
		helpers.ErrorResponse(response, "Error while decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	addedSpot, err := spotRepo.AddSpot(request.Context(), spot)
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, addedSpot)
}

func getSpotById(response http.ResponseWriter, request *http.Request, id string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	spot, err := spotRepo.FindSpotById(request.Context(), id)
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, spot)
}

func updateSpotById(response http.ResponseWriter, request *http.Request, id string) {
	found, err := spotRepo.FindSpotById(request.Context(), id)
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	spot := found[0]

	if err := helpers.CanEditAsset(request, spot.AddedBy); err != nil {
		helpers.ErrorResponse(response, "Authorization error: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if err := helpers.DecodeAndValidateRequestBody(request, &spot); err != nil {
		helpers.ErrorResponse(response, "Error while decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedSpot, err := spotRepo.UpdateSpot(request.Context(), id, spot)
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, updatedSpot)
}

func deleteSpotById(response http.ResponseWriter, request *http.Request, id string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	found, err := spotRepo.FindSpotById(request.Context(), id)
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	spot := found[0]

	if err := helpers.CanEditAsset(request, spot.AddedBy); err != nil {
		helpers.ErrorResponse(response, "Authorization error: "+err.Error(), http.StatusUnauthorized)
		return
	}

	err = spotRepo.DeleteSpotById(request.Context(), id)
	if err != nil {
		helpers.HandleRepoError(response, err)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
