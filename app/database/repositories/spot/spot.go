package spot

import (
	"context"
	"fmt"
	"scenic-spots-api/app/database"
	common "scenic-spots-api/app/database/repositories/common"
	reviewRepo "scenic-spots-api/app/database/repositories/review"
	"scenic-spots-api/internal/repoerrors"
	"scenic-spots-api/models"
	"scenic-spots-api/utils/calc"
	"scenic-spots-api/utils/generics"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

func buildSpotQuery(collectionRef *firestore.CollectionRef, params models.SpotQueryParams) (firestore.Query, error) {
	query := collectionRef.Query

	if params.Name != "" {
		query = query.Where("name", "==", strings.ToLower(params.Name))
	}
	if params.Latitude != "" || params.Longitude != "" || params.Radius != "" {
		if params.Latitude == "" || params.Longitude == "" || params.Radius == "" {
			return firestore.Query{}, fmt.Errorf("invalid parameter: latitude, longitude, and radius must all be provided together")
		}
		coordinates, err := calc.CoordinatesAfterRadius(params.Latitude, params.Longitude, params.Radius)
		if err != nil {
			return firestore.Query{}, err
		}
		query = query.Where("latitude", "<=", coordinates.MaxLat).
			Where("latitude", ">=", coordinates.MinLat).
			Where("longitude", "<=", coordinates.MaxLon).
			Where("longitude", ">=", coordinates.MinLon)
	}
	if params.Category != "" {
		query = query.Where("category", "==", strings.ToLower(params.Category))
	}

	return query, nil
}

func GetSpot(ctx context.Context, params models.SpotQueryParams) ([]models.Spot, error) {
	client := database.GetFirestoreClient()
	collectionRef := client.Collection(database.SpotCollectionName)

	query, err := buildSpotQuery(collectionRef, params)
	if err != nil {
		return []models.Spot{}, &repoerrors.InvalidQueryParameterError{
			Message: err.Error(),
		}
	}

	found, err := common.GetAllItems[*models.Spot](ctx, query)
	if err != nil {
		return []models.Spot{}, err
	}

	result := generics.DereferenceAll(found)

	return result, nil
}

func AddSpot(ctx context.Context, spotInfo models.NewSpot) ([]models.Spot, error) {
	if err := checkIfSpotAlreadyExists(ctx, spotInfo.Latitude, spotInfo.Longitude); err != nil {
		return []models.Spot{}, err
	}

	spot := models.Spot{
		Name:        strings.ToLower(spotInfo.Name),
		Description: spotInfo.Description,
		Latitude:    spotInfo.Latitude,
		Longitude:   spotInfo.Longitude,
		Category:    strings.ToLower(spotInfo.Category),
		Photos:      []string{},
		AddedBy:     "test user", /* TODO */
		CreatedAt:   time.Now(),
	}

	addedSpot, err := common.AddItem(ctx, database.SpotCollectionName, &spot)
	if err != nil {
		return []models.Spot{}, err
	}

	var result []models.Spot
	result = append(result, *addedSpot)

	return result, nil
}

func FindSpotById(ctx context.Context, id string) ([]models.Spot, error) {
	spot, err := common.FindItemById[*models.Spot](ctx, database.SpotCollectionName, id)
	if err != nil {
		return []models.Spot{}, err
	}

	var result []models.Spot
	result = append(result, *spot)

	return result, nil
}

func UpdateSpot(ctx context.Context, id string, newValues models.NewSpot) ([]models.Spot, error) {
	var err error
	result := []models.Spot{}

	spotToUpdate, err := common.FindItemById[*models.Spot](ctx, database.SpotCollectionName, id)
	if err != nil {
		return []models.Spot{}, err
	}

	if newValues.Name != "" {
		spotToUpdate.Name = newValues.Name
	}
	if newValues.Description != "" {
		spotToUpdate.Description = newValues.Description
	}
	if newValues.Latitude != 0 && newValues.Longitude != 0 {
		if err := checkIfSpotAlreadyExists(ctx, newValues.Latitude, newValues.Longitude); err != nil {
			return []models.Spot{}, err
		}
	}
	if newValues.Latitude != 0 {
		spotToUpdate.Latitude = newValues.Latitude
	}
	if newValues.Longitude != 0 {
		spotToUpdate.Longitude = newValues.Longitude
	}
	if newValues.Category != "" {
		spotToUpdate.Category = newValues.Category
	}

	client := database.GetFirestoreClient()
	_, err = client.Collection(database.SpotCollectionName).Doc(id).Update(ctx, []firestore.Update{
		{Path: "name", Value: spotToUpdate.Name},
		{Path: "description", Value: spotToUpdate.Description},
		{Path: "latitude", Value: spotToUpdate.Latitude},
		{Path: "longitude", Value: spotToUpdate.Longitude},
		{Path: "category", Value: spotToUpdate.Category},
	})
	if err != nil {
		return []models.Spot{}, err
	}

	result = append(result, *spotToUpdate)

	return result, nil
}

func DeleteSpotById(ctx context.Context, id string) error {
	if _, err := common.FindItemById[*models.Spot](ctx, database.SpotCollectionName, id); err != nil {
		return err
	}

	if err := reviewRepo.DeleteAllReviews(ctx, id); err != nil {
		return err
	}

	return common.DeleteItemById(ctx, database.SpotCollectionName, id)
}

// Checking if any spot in 100meter radius exists!
func checkIfSpotAlreadyExists(ctx context.Context, latitude float64, longitude float64) error {
	client := database.GetFirestoreClient()
	collectionRef := client.Collection(database.SpotCollectionName)

	query, _ := buildSpotQuery(collectionRef, models.SpotQueryParams{
		Name:      "",
		Latitude:  strconv.FormatFloat(latitude, 'f', -1, 64),
		Longitude: strconv.FormatFloat(longitude, 'f', -1, 64),
		Radius:    "0.1",
	})

	docs := query.Documents(ctx)
	results, _ := docs.GetAll()
	if len(results) != 0 {
		return repoerrors.ErrSpotAlreadyExists
	}

	return nil
}
