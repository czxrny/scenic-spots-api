package helpers

import (
	"encoding/json"
	"errors"
	"net/http"
	"scenic-spots-api/internal/database/repositories/repoerrors"
)

func WriteJSONResponse(response http.ResponseWriter, status int, data any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	if err := json.NewEncoder(response).Encode(data); err != nil {
		ErrorResponse(response, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
	}
}

func HandleRepoError(response http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repoerrors.ErrDoesNotExist):
		ErrorResponse(response, "Item does not exist", http.StatusNotFound)
	case errors.Is(err, repoerrors.ErrAlreadyExists):
		ErrorResponse(response, "Conflict: resource already exists", http.StatusConflict)
	default:
		ErrorResponse(response, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
	}
}

func RequireAuthAndPermission(response http.ResponseWriter, request *http.Request, addedBy string) bool {
	if err := IsAuthenticated(request); err != nil {
		ErrorResponse(response, err.Error(), http.StatusUnauthorized)
		return false
	}
	if err := CanEditAsset(request, addedBy); err != nil {
		ErrorResponse(response, "Authorization error: "+err.Error(), http.StatusUnauthorized)
		return false
	}
	return true
}
