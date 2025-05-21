package apihelpers

import (
	"encoding/json"
	"net/http"
	"scenic-spots-api/models"
	"scenic-spots-api/utils/logger"
)

func ErrorResponse(response http.ResponseWriter, message string, statusCode int) {
	errorResponse := models.APIError{
		Code:    statusCode,
		Message: message,
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)

	if err := json.NewEncoder(response).Encode(errorResponse); err != nil {
		logger.Error("Failed to encode JSON error response: " + err.Error())
	}
}
