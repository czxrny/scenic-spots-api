package spot

import (
	"context"
	"net/url"
	"scenic-spots-api/internal/auth"
	reviewRepo "scenic-spots-api/internal/database/repositories/review"
	spotRepo "scenic-spots-api/internal/database/repositories/spot"
	"scenic-spots-api/internal/models"
	"time"
)

func GetSpot(ctx context.Context, query url.Values) ([]models.Spot, error) {
	params := models.SpotQueryParams{
		Name:      query.Get("name"),
		Latitude:  query.Get("latitude"),
		Longitude: query.Get("longitude"),
		Radius:    query.Get("radius"),
		Category:  query.Get("category"),
	}

	found, err := spotRepo.GetSpot(ctx, params)
	if err != nil {
		return []models.Spot{}, err
	}
	return found, nil
}

func AddSpot(ctx context.Context, token string, newSpotInfo models.NewSpot) ([]models.Spot, error) {
	userName, err := auth.ExtractFromToken(token, "usr")
	if err != nil {
		return []models.Spot{}, err
	}
	if err := spotRepo.CheckIfSpotAlreadyExists(ctx, newSpotInfo.Latitude, newSpotInfo.Longitude); err != nil {
		return []models.Spot{}, err
	}

	spot := models.Spot{
		Name:        newSpotInfo.Name,
		Description: newSpotInfo.Description,
		Latitude:    newSpotInfo.Latitude,
		Longitude:   newSpotInfo.Longitude,
		Category:    newSpotInfo.Category,
		Photos:      []string{},
		AddedBy:     userName,
		CreatedAt:   time.Now(),
	}

	addedSpot, err := spotRepo.AddSpot(ctx, spot)
	if err != nil {
		return []models.Spot{}, err
	}

	var result []models.Spot
	result = append(result, addedSpot)

	return result, nil
}

func FindSpotById(ctx context.Context, id string) ([]models.Spot, error) {
	spot, err := spotRepo.FindSpotById(ctx, id)
	if err != nil {
		return []models.Spot{}, err
	}

	var result []models.Spot
	result = append(result, spot)

	return result, nil
}

func UpdateSpotById(ctx context.Context, token string, newSpotInfo models.NewSpot, id string) ([]models.Spot, error) {
	if err := spotRepo.CheckIfSpotAlreadyExists(ctx, newSpotInfo.Latitude, newSpotInfo.Longitude); err != nil {
		return []models.Spot{}, err
	}

	spot, err := spotRepo.FindSpotById(ctx, id)
	if err != nil {
		return []models.Spot{}, err
	}

	if err := auth.IsAuthorizedToEditAsset(token, spot.AddedBy); err != nil {
		return []models.Spot{}, err
	}

	if err := spotRepo.UpdateSpot(ctx, id, newSpotInfo); err != nil {
		return []models.Spot{}, err
	}

	spot.Name = newSpotInfo.Name
	spot.Description = newSpotInfo.Description
	spot.Latitude = newSpotInfo.Latitude
	spot.Longitude = newSpotInfo.Longitude
	spot.Category = newSpotInfo.Category

	var result []models.Spot
	result = append(result, spot)

	return result, nil
}

func DeleteSpotById(ctx context.Context, token string, id string) error {
	spot, err := spotRepo.FindSpotById(ctx, id)
	if err != nil {
		return err
	}

	if err := auth.IsAuthorizedToEditAsset(token, spot.AddedBy); err != nil {
		return err
	}

	if err := reviewRepo.DeleteAllReviews(ctx, id); err != nil {
		return err
	}

	return spotRepo.DeleteSpotById(ctx, id)
}
