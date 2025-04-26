package handlers

import (
	"encoding/json"
	"net/http"
	"scenic-spots-api/app/logger"
	"scenic-spots-api/models"
)

func ErrorResponse(response http.ResponseWriter, code string, message string, statusCode int) {
	errorResponse := models.APIError{
		Code:    code,
		Message: message,
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)

	if err := json.NewEncoder(response).Encode(errorResponse); err != nil {
		logger.Error("JSON Encoding error")
	}
}
