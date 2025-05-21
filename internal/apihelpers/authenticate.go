package apihelpers

import (
	"fmt"
	"net/http"
	"scenic-spots-api/app/auth"
	"strings"
)

func IsAuthenticated(request *http.Request) error {
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return fmt.Errorf("missing or invalid Authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	return auth.VerifyToken(token)
}
