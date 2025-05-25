package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"scenic-spots-api/internal/api"
	"strings"
)

func DecodeSegment(seg string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(seg)
}

// User is able to edit asset if he was the one posting, or has admin role.
func IsAuthorizedToEditAsset(token string, originalUser string) error {
	user, err := ExtractFromToken(token, "usr")
	if err != nil {
		return err
	}

	role, err := ExtractFromToken(token, "rol")
	if err != nil {
		return err
	}

	if user != originalUser && role != "admin" {
		return api.ErrIsUnauthorized
	}

	return nil
}

func ExtractFromToken(token string, field string) (string, error) {
	claims, err := extractClaimsFromToken(token)
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

	payloadData, err := DecodeSegment(parts[1])
	if err != nil {
		return nil, fmt.Errorf("error decoding JWT payload: %w", err)
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payloadData, &claims); err != nil {
		return nil, fmt.Errorf("invalid JWT payload: %w", err)
	}

	return claims, nil
}
