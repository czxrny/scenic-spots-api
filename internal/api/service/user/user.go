package user

import (
	"context"
	"scenic-spots-api/internal/auth"
	"scenic-spots-api/internal/database/repositories/repoerrors"
	userAuthRepo "scenic-spots-api/internal/database/repositories/user"
	"scenic-spots-api/internal/models"
)

func RegisterUser(ctx context.Context, userRegisterInfo models.UserRegisterInfo) (models.UserTokenResponse, error) {
	if err := ensureCredentialsUniqueness(ctx, userRegisterInfo.Name, userRegisterInfo.Email); err != nil {
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
	user, err := userAuthRepo.GetUserByField(ctx, "email", credentials.Email)
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
	user, err := userAuthRepo.FindUserById(ctx, userId)
	if err != nil {
		return err
	}

	if err := auth.IsAuthorizedToEditAsset(token, user.Name); err != nil {
		return err
	}

	return userAuthRepo.DeleteUserById(ctx, userId)
}

func ensureCredentialsUniqueness(ctx context.Context, userName string, email string) error {
	if _, err := userAuthRepo.GetUserByField(ctx, "email", email); err != nil && err != repoerrors.ErrDoesNotExist {
		return err
	} else if err == nil {
		return repoerrors.ErrAlreadyExists
	}
	if _, err := userAuthRepo.GetUserByField(ctx, "name", userName); err != nil && err != repoerrors.ErrDoesNotExist {
		return err
	} else if err == nil {
		return repoerrors.ErrAlreadyExists
	}
	return nil
}
