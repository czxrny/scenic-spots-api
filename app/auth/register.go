package auth

import (
	"context"
	"fmt"
	userRepo "scenic-spots-api/app/database/repositories/user"
	"scenic-spots-api/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx context.Context, credentials models.UserCredentials) (models.UserTokenResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserTokenResponse{}, fmt.Errorf("Error hashing password")
	}
	credentials.Password = string(hashed)

	addedUser, err := userRepo.AddUser(ctx, credentials)
	if err != nil {
		return models.UserTokenResponse{}, err
	}

	token, err := CreateToken(addedUser[0])
	if err != nil {
		return models.UserTokenResponse{}, err
	}

	return models.UserTokenResponse{
		Token:   token,
		LocalId: addedUser[0].Id,
	}, nil
}
