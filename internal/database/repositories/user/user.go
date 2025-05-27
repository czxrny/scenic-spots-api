package user

import (
	"context"
	"scenic-spots-api/internal/database"
	"scenic-spots-api/internal/database/repositories/common"
	"scenic-spots-api/internal/database/repositories/repoerrors"
	"scenic-spots-api/internal/models"
)

func AddUser(ctx context.Context, newUser models.User) (models.User, error) {
	addedUser, err := common.AddItem(ctx, models.UserAuthCollectionName, &newUser)
	if err != nil {
		return models.User{}, err
	}

	return *addedUser, nil
}

func FindUserById(ctx context.Context, id string) (models.User, error) {
	spot, err := common.FindItemById[*models.User](ctx, models.UserAuthCollectionName, id)
	if err != nil {
		return models.User{}, err
	}

	return *spot, nil
}

func DeleteUserById(ctx context.Context, id string) error {
	return common.DeleteItemById(ctx, models.UserAuthCollectionName, id)
}

func GetUserByField(ctx context.Context, fieldName, value string) (*models.User, error) {
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
