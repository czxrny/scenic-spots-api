package handlers

import (
	"io"
	"net/http"
)

func requestBodyIsEmpty(request *http.Request) bool {
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return false
	}
	if len(bodyBytes) > 0 {
		return false
	}

	return true
}
