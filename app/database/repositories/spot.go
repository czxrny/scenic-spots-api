package repositories

import (
	"context"
	"fmt"
	"scenic-spots-api/app/database"
	"scenic-spots-api/app/logger"
	"scenic-spots-api/models"
	"scenic-spots-api/utils/calc"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func buildQuery(collectionRef *firestore.CollectionRef, params models.SpotQueryParams) (*firestore.Query, error) {
	query := collectionRef.Query

	if params.Name != "" {
		query = query.Where("name", "==", params.Name)
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
		query = query.Where("category", "==", params.Category)
	}

	return &query, nil
}

func FindSpotsByQueryParams(params models.SpotQueryParams, ctx context.Context) (models.SpotMap, error) {
	client := database.GetFirestoreClient()
	collectionRef := client.Collection("spots")

	query, err := buildQuery(collectionRef, params)
	if err != nil {
		return models.SpotMap{}, err
	}
	iter := query.Documents(ctx)

	found := models.SpotMap{
		Spots: make(map[string]models.Spot),
	}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error("errror iteratora?")
			return models.SpotMap{}, err
		}

		var spot models.Spot

		if err := doc.DataTo(&spot); err != nil {
			logger.Error(err.Error())
			return models.SpotMap{}, err
		}

		found.Spots[doc.Ref.ID] = models.Spot{
			Latitude:  spot.Latitude,
			Longitude: spot.Longitude,
			Category:  spot.Category,
			Photos:    spot.Photos,
			AddedBy:   spot.AddedBy,
			CreatedAt: spot.CreatedAt,
		}
	}
	return found, nil
}
