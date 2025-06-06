package auth

import (
	"context"
	"scenic-spots-api/internal/api/apierrors"
	"scenic-spots-api/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func ValidatePasswordAndReturnToken(ctx context.Context, userInfo models.User, password string) (string, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(password)); err != nil {
		return "", apierrors.ErrInvalidCredentials
	}

	token, err := CreateToken(userInfo)
	if err != nil {
		return "", err
	}

	return token, nil
}
