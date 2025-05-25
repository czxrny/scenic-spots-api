package auth

import (
	"context"
	"fmt"
	"scenic-spots-api/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func ValidatePasswordAndReturnToken(ctx context.Context, userInfo models.User, password string) (string, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(userInfo.Password)); err != nil {
		return "", fmt.Errorf("Wrong password.")
	}

	token, err := CreateToken(userInfo)
	if err != nil {
		return "", err
	}

	return token, nil
}
