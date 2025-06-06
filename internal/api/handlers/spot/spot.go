package spot

import (
	"net/http"
	pHandler "scenic-spots-api/internal/api/handlers/photo"
	rHandler "scenic-spots-api/internal/api/handlers/review"
	helpers "scenic-spots-api/internal/api/helpers"
	spotService "scenic-spots-api/internal/api/service/spot"
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

	found, err := spotService.GetSpot(request.Context(), request.URL.Query())
	if err != nil {
		helpers.HandleErrors(response, err)
		return
	}
	helpers.WriteJSONResponse(response, http.StatusOK, found)
}

func addSpot(response http.ResponseWriter, request *http.Request) {
	token, err := helpers.GetJWTToken(request)
	if err != nil {
		helpers.ErrorResponse(response, "Error while decoding header: "+err.Error(), http.StatusBadRequest)
		return
	}

	var spot models.NewSpot
	if err := helpers.DecodeAndValidateRequestBody(request, &spot); err != nil {
		helpers.ErrorResponse(response, "Error while decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := spotService.AddSpot(request.Context(), token, spot)
	if err != nil {
		helpers.HandleErrors(response, err)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, result)
}

func getSpotById(response http.ResponseWriter, request *http.Request, id string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}

	result, err := spotService.FindSpotById(request.Context(), id)
	if err != nil {
		helpers.HandleErrors(response, err)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, result)
}

func updateSpotById(response http.ResponseWriter, request *http.Request, id string) {
	token, err := helpers.GetJWTToken(request)
	if err != nil {
		helpers.ErrorResponse(response, "Error while decoding header: "+err.Error(), http.StatusBadRequest)
		return
	}

	var spot models.NewSpot
	if err := helpers.DecodeAndValidateRequestBody(request, &spot); err != nil {
		helpers.ErrorResponse(response, "Error while decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := spotService.UpdateSpotById(request.Context(), token, spot, id)
	if err != nil {
		helpers.HandleErrors(response, err)
		return
	}

	helpers.WriteJSONResponse(response, http.StatusOK, result)
}

func deleteSpotById(response http.ResponseWriter, request *http.Request, id string) {
	if !helpers.RequestBodyIsEmpty(request) {
		helpers.ErrorResponse(response, "GET request must not contain a body", http.StatusBadRequest)
		return
	}
	token, err := helpers.GetJWTToken(request)
	if err != nil {
		helpers.ErrorResponse(response, "Error while decoding header: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = spotService.DeleteSpotById(request.Context(), token, id)
	if err != nil {
		helpers.HandleErrors(response, err)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
