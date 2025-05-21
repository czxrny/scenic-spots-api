package auth

import (
	"context"
	"fmt"
	userRepo "scenic-spots-api/internal/database/repositories/user"
	"scenic-spots-api/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx context.Context, userRegisterInfo models.UserRegisterInfo) (models.UserTokenResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(userRegisterInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserTokenResponse{}, fmt.Errorf("Error hashing password")
	}
	userRegisterInfo.Password = string(hashed)

	addedUser, err := userRepo.AddUser(ctx, userRegisterInfo)
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
