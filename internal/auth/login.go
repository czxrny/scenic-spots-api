package auth

import (
	"context"
	"fmt"
	userRepo "scenic-spots-api/internal/database/repositories/user"
	"scenic-spots-api/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func LoginUser(ctx context.Context, credentials models.UserCredentials) (models.UserTokenResponse, error) {
	user, err := userRepo.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		return models.UserTokenResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return models.UserTokenResponse{}, fmt.Errorf("Wrong password.")
	}

	token, err := CreateToken(*user)
	if err != nil {
		return models.UserTokenResponse{}, err
	}

	return models.UserTokenResponse{
		Token:   token,
		LocalId: user.Id,
	}, nil
}
