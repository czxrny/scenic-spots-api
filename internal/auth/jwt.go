package auth

import (
	"encoding/base64"
	"fmt"
	"os"
	"scenic-spots-api/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(time.Hour)

	claims := jwt.MapClaims{
		"lid": user.Id,
		"usr": user.Name,
		"rol": user.Role,
		"exp": expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("token has a wrong signature algorithm set")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return fmt.Errorf("Bad token!: " + err.Error())
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return fmt.Errorf("Invalid token formatting / expired", err)
	}

	return nil
}

func DecodeSegment(seg string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(seg)
}
