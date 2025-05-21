package apihelpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"scenic-spots-api/app/auth"
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

func CanEditAsset(request *http.Request, orignalId string) error {
	token, err := GetJWTToken(request)
	if err != nil {
		return err
	}

	userId, err := extractFromToken(&token, "lid")
	if err != nil {
		return err
	}

	role, err := extractFromToken(&token, "rol")
	if err != nil {
		return err
	}

	if userId != orignalId && role != "admin" {
		return fmt.Errorf("illegal operation: not authorized to delete this user")
	}

	return nil
}

func extractFromToken(token *string, field string) (string, error) {
	claims, err := extractClaimsFromToken(*token)
	if err != nil {
		return "", err
	}

	fieldVal, ok := claims[field].(string)
	if !ok {
		return "", fmt.Errorf("%v not found or is not a string", field)
	}
	return fieldVal, nil
}

func extractClaimsFromToken(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT format")
	}

	payloadData, err := auth.DecodeSegment(parts[1])
	if err != nil {
		return nil, fmt.Errorf("error decoding JWT payload: %w", err)
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payloadData, &claims); err != nil {
		return nil, fmt.Errorf("invalid JWT payload: %w", err)
	}

	return claims, nil
}
