package spot

import (
	"context"
	"net/url"
	"scenic-spots-api/internal/api/apierrors"
	"scenic-spots-api/internal/auth"
	"scenic-spots-api/internal/database/repositories/repoerrors"
	reviewRepo "scenic-spots-api/internal/database/repositories/review"
	spotRepo "scenic-spots-api/internal/database/repositories/spot"
	"scenic-spots-api/internal/models"
	"strconv"
	"time"
)

func GetSpot(ctx context.Context, query url.Values) ([]models.Spot, error) {
	params := models.SpotQueryParams{
		Name:      query.Get("name"),
		Latitude:  query.Get("latitude"),
		Longitude: query.Get("longitude"),
		Radius:    query.Get("radius"),
		Category:  query.Get("category"),
		AddedBy:   query.Get("addedBy"),
	}

	if (params.Latitude != "" || params.Longitude != "" || params.Radius != "") &&
		(params.Latitude == "" || params.Longitude == "" || params.Radius == "") {
		return nil, apierrors.ErrInvalidQueryParameters
	}

	spots, err := spotRepo.GetSpot(ctx, params)
	if err != nil {
		return nil, err
	}
	return spots, nil
}

func AddSpot(ctx context.Context, token string, newSpotInfo models.NewSpot) (models.Spot, error) {
	userName, err := auth.ExtractFromToken(token, "usr")
	if err != nil {
		return models.Spot{}, err
	}

	if found, err := getNearbySpots(ctx, newSpotInfo.Latitude, newSpotInfo.Longitude); err != nil {
		return models.Spot{}, err
	} else if len(found) > 0 {
		return models.Spot{}, repoerrors.ErrAlreadyExists
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
		return models.Spot{}, err
	}

	return addedSpot, nil
}

func FindSpotById(ctx context.Context, id string) (models.Spot, error) {
	spot, err := spotRepo.FindSpotById(ctx, id)
	if err != nil {
		return models.Spot{}, err
	}

	return spot, nil
}

func UpdateSpotById(ctx context.Context, token string, newSpotInfo models.NewSpot, id string) (models.Spot, error) {
	if found, err := getNearbySpots(ctx, newSpotInfo.Latitude, newSpotInfo.Longitude); err != nil {
		return models.Spot{}, err
	} else if len(found) > 0 && found[0].Id != id {
		return models.Spot{}, repoerrors.ErrAlreadyExists
	}

	spot, err := spotRepo.FindSpotById(ctx, id)
	if err != nil {
		return models.Spot{}, err
	}

	if err := auth.IsAuthorizedToEditAsset(token, spot.AddedBy); err != nil {
		return models.Spot{}, err
	}

	if err := spotRepo.UpdateSpot(ctx, id, newSpotInfo); err != nil {
		return models.Spot{}, err
	}

	spot.Name = newSpotInfo.Name
	spot.Description = newSpotInfo.Description
	spot.Latitude = newSpotInfo.Latitude
	spot.Longitude = newSpotInfo.Longitude
	spot.Category = newSpotInfo.Category

	return spot, nil
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

// Check if there are no spots in 100m radius!
func getNearbySpots(ctx context.Context, latitude float64, longitude float64) ([]models.Spot, error) {
	found, err := spotRepo.GetSpot(ctx, models.SpotQueryParams{
		Name:      "",
		Latitude:  strconv.FormatFloat(latitude, 'f', -1, 64),
		Longitude: strconv.FormatFloat(longitude, 'f', -1, 64),
		Radius:    "0.1",
		Category:  "",
	})
	if err != nil {
		return []models.Spot{}, err
	}
	return found, nil
}
