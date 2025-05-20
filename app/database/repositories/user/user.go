package user

import (
	"context"
	"fmt"
	"scenic-spots-api/app/database"
	"scenic-spots-api/app/database/repositories/common"
	"scenic-spots-api/internal/repoerrors"
	"scenic-spots-api/models"

	"golang.org/x/crypto/bcrypt"
)

func AddUser(ctx context.Context, credentials models.UserCredentials) ([]models.User, error) {
	if err := checkIfEmailAlreadyExists(ctx, credentials.Email); err != nil {
		return []models.User{}, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return []models.User{}, fmt.Errorf("Error hashing password")
	}

	newUser := models.User{
		Name:     credentials.Name,
		Email:    credentials.Email,
		Password: string(hashed),
		Role:     "user",
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
	client := database.GetFirestoreClient()
	result, _ := client.Collection(database.UserAuthCollectionName).Where("email", "==", email).Documents(ctx).GetAll()

	if len(result) != 0 {
		return repoerrors.ErrAlreadyExists
	}
	return nil
}
