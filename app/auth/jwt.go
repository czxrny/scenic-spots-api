package auth

import (
	"os"
	"scenic-spots-api/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"lid":  user.Id,
		"usr":  user.Name,
		"role": user.Role,
		"exp":  expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
