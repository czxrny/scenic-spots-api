package user

import (
	"context"
	"scenic-spots-api/internal/auth"
	userAuthRepo "scenic-spots-api/internal/database/repositories/user"
	"scenic-spots-api/internal/models"
)

func RegisterUser(ctx context.Context, userRegisterInfo models.UserRegisterInfo) (models.UserTokenResponse, error) {
	if _, err := userAuthRepo.CheckIfEmailExists(ctx, userRegisterInfo.Email); err != nil {
		return models.UserTokenResponse{}, err
	}
	if _, err := userAuthRepo.CheckIfUsernameExists(ctx, userRegisterInfo.Name); err != nil {
		return models.UserTokenResponse{}, err
	}

	if err := auth.EncryptThePassword(&userRegisterInfo); err != nil {
		return models.UserTokenResponse{}, err
	}

	newUser := models.User{
		Name:     userRegisterInfo.Name,
		Email:    userRegisterInfo.Email,
		Password: userRegisterInfo.Password,
		Role:     "user", // by default
	}

	addedUser, err := userAuthRepo.AddUser(ctx, newUser)
	if err != nil {
		return models.UserTokenResponse{}, err
	}

	token, err := auth.CreateToken(addedUser)
	if err != nil {
		return models.UserTokenResponse{}, err
	}

	return models.UserTokenResponse{
		Token:   token,
		LocalId: addedUser.Id,
	}, nil
}

func LoginUser(ctx context.Context, credentials models.UserCredentials) (models.UserTokenResponse, error) {
	user, err := userAuthRepo.CheckIfEmailExists(ctx, credentials.Email)
	if err != nil {
		return models.UserTokenResponse{}, err
	}

	token, err := auth.ValidatePasswordAndReturnToken(ctx, *user, credentials.Password)
	if err != nil {
		return models.UserTokenResponse{}, err
	}

	return models.UserTokenResponse{
		Token:   token,
		LocalId: user.Id,
	}, nil
}

func DeleteUserById(ctx context.Context, token string, userId string) error {
	_, err := userAuthRepo.FindUserById(ctx, userId)
	if err != nil {
		return err
	}

	if err := auth.IsAuthorizedToEditAsset(token, userId); err != nil {
		return err
	}

	return userAuthRepo.DeleteUserById(ctx, userId)
}
