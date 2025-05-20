package user

import (
	"context"
	"scenic-spots-api/app/database"
	"scenic-spots-api/app/database/repositories/common"
	"scenic-spots-api/internal/repoerrors"
	"scenic-spots-api/models"
)

func AddUser(ctx context.Context, credentials models.UserCredentials) ([]models.User, error) {
	if err := checkIfEmailAlreadyExists(ctx, credentials.Email); err != nil {
		return []models.User{}, err
	}

	newUser := models.User{
		Name:     credentials.Name,
		Email:    credentials.Email,
		Password: credentials.Password,
		Role:     "user", // by default
	}

	addedUser, err := common.AddItem(ctx, database.UserAuthCollectionName, &newUser)
	if err != nil {
		return []models.User{}, err
	}

	var result []models.User
	result = append(result, *addedUser)

	return result, nil
}

func checkIfEmailAlreadyExists(ctx context.Context, email string) error {
	result, err := GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if result != nil {
		return repoerrors.ErrAlreadyExists
	}
	return nil
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	client := database.GetFirestoreClient()
	query := client.Collection(database.UserAuthCollectionName).Where("email", "==", email)

	result, err := common.GetAllItems[*models.User](ctx, query)
	if err != nil {
		return nil, err
	}

	return result[0], nil
}
