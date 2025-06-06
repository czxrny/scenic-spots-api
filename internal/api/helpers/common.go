package helpers

import (
	"encoding/json"
	"errors"
	"net/http"
	"scenic-spots-api/internal/api/apierrors"
	"scenic-spots-api/internal/database/repositories/repoerrors"
)

func WriteJSONResponse(response http.ResponseWriter, status int, data any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	if err := json.NewEncoder(response).Encode(data); err != nil {
		ErrorResponse(response, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
	}
}

func HandleErrors(response http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repoerrors.ErrDoesNotExist):
		ErrorResponse(response, "Item does not exist", http.StatusNotFound)
	case errors.Is(err, repoerrors.ErrAlreadyExists):
		ErrorResponse(response, "Conflict: resource already exists", http.StatusConflict)
	case errors.Is(err, apierrors.ErrInvalidQueryParameters):
		ErrorResponse(response, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
	case errors.Is(err, apierrors.ErrInvalidCredentials):
		ErrorResponse(response, "Authorization error: "+err.Error(), http.StatusUnauthorized)
	case errors.Is(err, apierrors.ErrIsUnauthorized):
		ErrorResponse(response, "Permission error: "+err.Error(), http.StatusForbidden)
	default:
		ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
	}
}
