package repositories

import (
	"context"
	"fmt"
	"scenic-spots-api/app/database"
	"scenic-spots-api/app/logger"
	"scenic-spots-api/models"
	"scenic-spots-api/utils/calc"
	"scenic-spots-api/utils/generics"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

const spotCollectionName string = "spots"

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
			logger.Error(err.Error())
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

func GetSpot(params models.SpotQueryParams, ctx context.Context) ([]models.Spot, error) {
	client := database.GetFirestoreClient()
	collectionRef := client.Collection(spotCollectionName)

	query, err := buildSpotQuery(collectionRef, params)
	if err != nil {
		return []models.Spot{}, err
	}

	found, err := getAllItems[*models.Spot](ctx, query)
	if err != nil {
		return []models.Spot{}, err
	}

	result := generics.DereferenceAll(found)

	return result, nil
}

func AddSpot(spotInfo models.NewSpot, ctx context.Context) ([]models.Spot, error) {
	client := database.GetFirestoreClient()
	collectionRef := client.Collection(spotCollectionName)

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

	// Casting to a json to avoid capitalized words in database.
	spotData := map[string]interface{}{
		"name":        spot.Name,
		"description": spot.Description,
		"latitude":    spot.Latitude,
		"longitude":   spot.Longitude,
		"category":    spot.Category,
		"photos":      spot.Photos,
		"addedBy":     spot.AddedBy,
		"createdAt":   spot.CreatedAt,
	}

	docRef, _, err := collectionRef.Add(ctx, spotData)

	if err != nil {
		return []models.Spot{}, err
	}

	spot.Id = docRef.ID
	var result []models.Spot

	result = append(result, spot)

	return result, nil
}

func FindSpotById(ctx context.Context, id string) ([]models.Spot, error) {
	spot, err := findItemById[*models.Spot](ctx, spotCollectionName, id)
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

	spotToUpdate, err := findItemById[*models.Spot](ctx, spotCollectionName, id)
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
	_, err = client.Collection(spotCollectionName).Doc(id).Update(ctx, []firestore.Update{
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
	return deleteItemById(ctx, spotCollectionName, id)
}

// Checking if any spot in 100meter radius exists!
func checkIfSpotAlreadyExists(ctx context.Context, latitude float64, longitude float64) error {
	client := database.GetFirestoreClient()
	collectionRef := client.Collection(spotCollectionName)

	query, err := buildSpotQuery(collectionRef, models.SpotQueryParams{
		Name:      "",
		Latitude:  strconv.FormatFloat(latitude, 'f', -1, 64),
		Longitude: strconv.FormatFloat(longitude, 'f', -1, 64),
		Radius:    "0.1",
	})
	if err != nil {
		return err
	}

	docs := query.Documents(ctx)
	results, _ := docs.GetAll()
	if len(results) != 0 {
		return fmt.Errorf("The spot already exists in the database: %s", results[0].Ref.ID)
	}

	return nil
}
