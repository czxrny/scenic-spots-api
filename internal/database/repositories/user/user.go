package user

import (
	"context"
	"scenic-spots-api/internal/database"
	"scenic-spots-api/internal/database/repositories/common"
	"scenic-spots-api/internal/database/repositories/repoerrors"
	"scenic-spots-api/internal/models"
)

func AddUser(ctx context.Context, userRegisterInfo models.UserRegisterInfo) ([]models.User, error) {
	if err := checkIfEmailAlreadyExists(ctx, userRegisterInfo.Email); err != nil {
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

func checkIfEmailAlreadyExists(ctx context.Context, email string) error {
	result, err := GetUserByEmail(ctx, email)
	if err == repoerrors.ErrDoesNotExist {
		return nil
	} else if result != nil {
		return repoerrors.ErrAlreadyExists
	} else {
		return err
	}
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	client := database.GetFirestoreClient()
	query := client.Collection(models.UserAuthCollectionName).Where("email", "==", email)

	result, err := common.GetAllItems[*models.User](ctx, query)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, repoerrors.ErrDoesNotExist
	}

	return result[0], nil
}

func DeleteUserById(ctx context.Context, id string) error {
	if _, err := common.FindItemById[*models.User](ctx, models.UserAuthCollectionName, id); err != nil {
		return err
	}

	return common.DeleteItemById(ctx, models.UserAuthCollectionName, id)
}
