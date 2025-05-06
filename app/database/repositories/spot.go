package repositories

import (
	"context"
	"fmt"
	"scenic-spots-api/app/database"
	"scenic-spots-api/app/logger"
	"scenic-spots-api/models"
	"scenic-spots-api/utils/calc"
	"strconv"
	"time"

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
			return models.SpotMap{}, err
		}

		var spot models.Spot

		if err := doc.DataTo(&spot); err != nil {
			return models.SpotMap{}, err
		}

		found.Spots[doc.Ref.ID] = models.Spot{
			Name:      spot.Name,
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

func AddSpot(spotInfo models.NewSpot, ctx context.Context) (models.SpotMap, error) {
	client := database.GetFirestoreClient()
	collectionRef := client.Collection("spots")

	if err := checkIfAlreadyExists(ctx, collectionRef, strconv.FormatFloat(spotInfo.Latitude, 'f', -1, 64), strconv.FormatFloat(spotInfo.Longitude, 'f', -1, 64)); err != nil {
		return models.SpotMap{}, err
	}

	spot := models.Spot{
		Name:      spotInfo.Name,
		Latitude:  spotInfo.Latitude,
		Longitude: spotInfo.Longitude,
		Category:  spotInfo.Category,
		Photos:    []string{},
		AddedBy:   "test user", /* TODO */
		CreatedAt: time.Now(),
	}

	// Casting to a json to avoid capitalized words in database.
	spotData := map[string]interface{}{
		"name":      spot.Name,
		"latitude":  spot.Latitude,
		"longitude": spot.Longitude,
		"category":  spot.Category,
		"photos":    spot.Photos,
		"addedBy":   spot.AddedBy,
		"createdAt": spot.CreatedAt, // <-- nadal time.Time
	}

	docRef, _, err := collectionRef.Add(ctx, spotData)

	if err != nil {
		return models.SpotMap{}, err
	}

	result := models.SpotMap{
		Spots: make(map[string]models.Spot),
	}

	result.Spots[docRef.ID] = spot

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
