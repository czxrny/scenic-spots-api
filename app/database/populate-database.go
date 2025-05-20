package database

import (
	"context"
	"encoding/json"
	"os"
	"scenic-spots-api/models"
	"scenic-spots-api/utils/generics"

	"cloud.google.com/go/firestore"
)

func populateDatabase(ctx context.Context) error {
	err := addExampleData[models.Spot](ctx, SpotCollectionName, os.Getenv("DB_SPOTS"))
	if err != nil {
		return err
	}

	return addExampleData[models.Review](ctx, ReviewCollectionName, os.Getenv("DB_REVIEWS"))
}

func addExampleData[T any](ctx context.Context, collectionName string, filePath string) error {
	itemMap, err := readFileToStruct[T](filePath)
	if err != nil {
		return err
	}

	return addToDatabase(ctx, GetFirestoreClient(), collectionName, itemMap)
}

func readFileToStruct[T any](filePath string) (map[string]T, error) {
	rawData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var itemMap map[string]T
	err = json.Unmarshal(rawData, &itemMap)
	if err != nil {
		return nil, err
	}
	return itemMap, nil
}

func addToDatabase[T ~map[string]V, V any](ctx context.Context, client *firestore.Client, collectionName string, items T) error {
	for id, item := range items {
		jsonItem, err := generics.StructToMapLower(item)
		if err != nil {
			return err
		}
		_, err = client.Collection(collectionName).Doc(id).Set(ctx, jsonItem)
		if err != nil {
			return err
		}
	}
	return nil
}
