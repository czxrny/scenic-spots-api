package spot

import (
	"context"
	"scenic-spots-api/internal/api/apierrors"
	"scenic-spots-api/internal/database"
	common "scenic-spots-api/internal/database/repositories/common"
	"scenic-spots-api/internal/models"
	"scenic-spots-api/utils/calc"
	"scenic-spots-api/utils/generics"
	"strings"

	"cloud.google.com/go/firestore"
)

func buildSpotQuery(collectionRef *firestore.CollectionRef, params models.SpotQueryParams) (firestore.Query, error) {
	query := collectionRef.Query

	if params.Name != "" {
		query = query.Where("name", "==", strings.ToLower(params.Name))
	}

	if params.Latitude != "" {
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
	collectionRef := client.Collection(models.SpotCollectionName)

	query, err := buildSpotQuery(collectionRef, params)
	if err != nil {
		return []models.Spot{}, &apierrors.InvalidQueryParameterError{
			Message: err.Error(),
		}
	}

	found, err := common.GetAllItems[*models.Spot](ctx, query)
	if err != nil {
		return []models.Spot{}, err
	}

	return generics.DereferenceAll(found), nil
}

func AddSpot(ctx context.Context, spot models.Spot) (models.Spot, error) {
	addedSpot, err := common.AddItem(ctx, models.SpotCollectionName, &spot)
	if err != nil {
		return models.Spot{}, err
	}

	return *addedSpot, nil
}

func FindSpotById(ctx context.Context, id string) (models.Spot, error) {
	spot, err := common.FindItemById[*models.Spot](ctx, models.SpotCollectionName, id)
	if err != nil {
		return models.Spot{}, err
	}

	return *spot, nil
}

func UpdateSpot(ctx context.Context, id string, updatedSpot models.NewSpot) error {
	client := database.GetFirestoreClient()
	_, err := client.Collection(models.SpotCollectionName).Doc(id).Update(ctx, []firestore.Update{
		{Path: "name", Value: updatedSpot.Name},
		{Path: "description", Value: updatedSpot.Description},
		{Path: "latitude", Value: updatedSpot.Latitude},
		{Path: "longitude", Value: updatedSpot.Longitude},
		{Path: "category", Value: updatedSpot.Category},
	})
	return err
}

func DeleteSpotById(ctx context.Context, id string) error {
	if _, err := common.FindItemById[*models.Spot](ctx, models.SpotCollectionName, id); err != nil {
		return err
	}

	return common.DeleteItemById(ctx, models.SpotCollectionName, id)
}
