package apihelpers

import (
	"io"
	"net/http"
)

func RequestBodyIsEmpty(request *http.Request) bool {
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return false
	}
	if len(bodyBytes) > 0 {
		return false
	}

	return true
}
