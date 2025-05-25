package helpers

import (
	"fmt"
	"net/http"
	"scenic-spots-api/internal/auth"
	"strings"
)

func GetJWTToken(request *http.Request) (string, error) {
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("missing or invalid Authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}

func IsAuthenticated(request *http.Request) error {
	token, err := GetJWTToken(request)
	if err != nil {
		return err
	}

	return auth.VerifyToken(token)
}
