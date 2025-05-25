package user

import (
	"context"
	"fmt"
	"scenic-spots-api/internal/database"
	"scenic-spots-api/internal/database/repositories/common"
	"scenic-spots-api/internal/database/repositories/repoerrors"
	"scenic-spots-api/internal/models"
)

func AddUser(ctx context.Context, userRegisterInfo models.UserRegisterInfo) ([]models.User, error) {
	if _, err := CheckIfEmailExists(ctx, userRegisterInfo.Email); err != nil {
		return []models.User{}, err
	}
	if _, err := CheckIfUsernameExists(ctx, userRegisterInfo.Name); err != nil {
		return []models.User{}, err
	}

	newUser := models.User{
		Name:     userRegisterInfo.Name,
		Email:    userRegisterInfo.Email,
		Password: userRegisterInfo.Password,
		Role:     "user", // by default
	}

	addedUser, err := common.AddItem(ctx, models.UserAuthCollectionName, &newUser)
	if err != nil {
		return []models.User{}, err
	}

	var result []models.User
	result = append(result, *addedUser)

	return result, nil
}

func CheckIfEmailExists(ctx context.Context, email string) (*models.User, error) {
	user, err := getUserByField(ctx, "email", email)
	if err != nil && err != repoerrors.ErrDoesNotExist {
		return nil, fmt.Errorf("error checking email: %w", err)
	}
	if user != nil {
		return user, repoerrors.ErrEmailAlreadyExists
	}
	return nil, nil
}

func CheckIfUsernameExists(ctx context.Context, username string) (*models.User, error) {
	user, err := getUserByField(ctx, "name", username)
	if err != nil && err != repoerrors.ErrDoesNotExist {
		return nil, fmt.Errorf("error checking username: %w", err)
	}
	if user != nil {
		return user, repoerrors.ErrUsernameAlreadyExists
	}
	return nil, nil
}

func getUserByField(ctx context.Context, fieldName, value string) (*models.User, error) {
	client := database.GetFirestoreClient()
	query := client.Collection(models.UserAuthCollectionName).Where(fieldName, "==", value)

	results, err := common.GetAllItems[*models.User](ctx, query)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, repoerrors.ErrDoesNotExist
	}

	return results[0], nil
}

func DeleteUserById(ctx context.Context, id string) error {
	if _, err := common.FindItemById[*models.User](ctx, models.UserAuthCollectionName, id); err != nil {
		return err
	}

	return common.DeleteItemById(ctx, models.UserAuthCollectionName, id)
}
