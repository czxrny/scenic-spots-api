package repositories

import (
	"context"
	"fmt"
	"scenic-spots-api/app/database"
	"scenic-spots-api/app/logger"
	"scenic-spots-api/models"
	"scenic-spots-api/utils/calc"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

var collectionName string = "spots"

func buildQuery(collectionRef *firestore.CollectionRef, params models.SpotQueryParams) (*firestore.Query, error) {
	query := collectionRef.Query

	if params.Name != "" {
		query = query.Where("name", "==", strings.ToLower(params.Name))
	}
	if params.Latitude != "" || params.Longitude != "" || params.Radius != "" {
		if params.Latitude == "" || params.Longitude == "" || params.Radius == "" {
			return nil, fmt.Errorf("invalid parameter: latitude, longitude, and radius must all be provided together")
		}
		coordinates, err := calc.CoordinatesAfterRadius(params.Latitude, params.Longitude, params.Radius)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		query = query.Where("latitude", "<=", coordinates.MaxLat).
			Where("latitude", ">=", coordinates.MinLat).
			Where("longitude", "<=", coordinates.MaxLon).
			Where("longitude", ">=", coordinates.MinLon)
	}
	if params.Category != "" {
		query = query.Where("category", "==", strings.ToLower(params.Category))
	}

	return &query, nil
}

func FindSpotsByQueryParams(params models.SpotQueryParams, ctx context.Context) ([]models.Spot, error) {
	client := database.GetFirestoreClient()
	collectionRef := client.Collection(collectionName)

	query, err := buildQuery(collectionRef, params)
	if err != nil {
		return []models.Spot{}, err
	}
	iter := query.Documents(ctx)

	var found []models.Spot

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []models.Spot{}, err
		}

		var spot models.Spot

		if err := doc.DataTo(&spot); err != nil {
			return []models.Spot{}, err
		}

		found = append(found, models.Spot{
			Id:          doc.Ref.ID,
			Name:        spot.Name,
			Description: spot.Description,
			Latitude:    spot.Latitude,
			Longitude:   spot.Longitude,
			Category:    spot.Category,
			Photos:      spot.Photos,
			AddedBy:     spot.AddedBy,
			CreatedAt:   spot.CreatedAt,
		})
	}
	return found, nil
}

func AddSpot(spotInfo models.NewSpot, ctx context.Context) ([]models.Spot, error) {
	client := database.GetFirestoreClient()
	collectionRef := client.Collection(collectionName)

	if err := checkIfAlreadyExists(ctx, collectionRef, strconv.FormatFloat(spotInfo.Latitude, 'f', -1, 64), strconv.FormatFloat(spotInfo.Longitude, 'f', -1, 64)); err != nil {
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

// Checking if any spot in 100meter radius exists!
func checkIfAlreadyExists(ctx context.Context, collectionRef *firestore.CollectionRef, latitude string, longitude string) error {
	query, err := buildQuery(collectionRef, models.SpotQueryParams{
		Name:      "",
		Latitude:  latitude,
		Longitude: longitude,
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

func FindById(ctx context.Context, id string) ([]models.Spot, error) {
	client := database.GetFirestoreClient()
	doc, err := client.Collection(collectionName).Doc(id).Get(ctx)
	if err != nil {
		return []models.Spot{}, err
	}

	var spot models.Spot
	if err := doc.DataTo(&spot); err != nil {
		return []models.Spot{}, err
	}

	// Name should won't ever be empty - unless no doc was found.
	if spot.Name == "" {
		return []models.Spot{}, fmt.Errorf("Spot with ID: %v does not exist", id)
	}

	spot.Id = doc.Ref.ID
	var result []models.Spot

	result = append(result, spot)

	return result, nil
}

func UpdateSpot(ctx context.Context, id string, newValues models.NewSpot) ([]models.Spot, error) {
	var err error
	result := []models.Spot{}

	result, err = FindById(ctx, id)
	if err != nil {
		return []models.Spot{}, err
	}

	spotToUpdate := result[0]
	if newValues.Name != "" {
		spotToUpdate.Name = newValues.Name
	}
	if newValues.Description != "" {
		spotToUpdate.Description = newValues.Description
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
	_, err = client.Collection(collectionName).Doc(id).Update(ctx, []firestore.Update{
		{Path: "name", Value: spotToUpdate.Name},
		{Path: "description", Value: spotToUpdate.Description},
		{Path: "latitude", Value: spotToUpdate.Latitude},
		{Path: "longitude", Value: spotToUpdate.Longitude},
		{Path: "category", Value: spotToUpdate.Category},
	})
	if err != nil {
		return []models.Spot{}, err
	}

	result[0] = spotToUpdate

	return result, nil
}

func DeleteById(ctx context.Context, id string) error {
	client := database.GetFirestoreClient()
	docRef := client.Collection(collectionName).Doc(id)

	_, err := docRef.Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
