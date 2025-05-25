package review

import (
	"context"
	"fmt"
	"scenic-spots-api/internal/database"
	common "scenic-spots-api/internal/database/repositories/common"
	"scenic-spots-api/internal/database/repositories/repoerrors"
	"scenic-spots-api/internal/models"
	"scenic-spots-api/utils/generics"
	"strconv"

	"cloud.google.com/go/firestore"
)

func buildReviewQuery(collectionRef *firestore.CollectionRef, params models.ReviewQueryParams) (firestore.Query, error) {
	var limit int
	var err error
	query := collectionRef.Query

	query = query.Where("spotId", "==", params.SpotId)

	if params.Limit != "" {
		limit, err = strconv.Atoi(params.Limit)
		if err != nil {
			return firestore.Query{}, fmt.Errorf("invalid limit parameter")
		}
		query = query.Limit(limit)
	}
	return query, nil
}

func GetReviews(ctx context.Context, params models.ReviewQueryParams) ([]models.Review, error) {
	if _, err := common.FindItemById[*models.Spot](ctx, models.SpotCollectionName, params.SpotId); err != nil {
		return []models.Review{}, err
	}

	client := database.GetFirestoreClient()
	collectionRef := client.Collection(models.ReviewCollectionName)

	query, err := buildReviewQuery(collectionRef, params)
	if err != nil {
		return []models.Review{}, &repoerrors.InvalidQueryParameterError{
			Message: err.Error(),
		}
	}

	found, err := common.GetAllItems[*models.Review](ctx, query)
	if err != nil {
		return []models.Review{}, err
	}

	result := generics.DereferenceAll(found)
	return result, nil
}

func AddReview(ctx context.Context, review models.Review) (models.Review, error) {
	addedReview, err := common.AddItem(ctx, models.SpotCollectionName, &review)
	if err != nil {
		return models.Review{}, err
	}

	return *addedReview, nil
}

func DeleteAllReviews(ctx context.Context, spotId string) error {
	client := database.GetFirestoreClient()
	query := client.Collection(models.ReviewCollectionName).Where("spotId", "==", spotId)

	return common.DeleteAllItems(ctx, query)
}

func FindReviewById(ctx context.Context, id string) (models.Review, error) {
	result, err := common.FindItemById[*models.Review](ctx, models.ReviewCollectionName, id)
	if err != nil {
		return models.Review{}, err
	}

	return *result, nil
}

func UpdateReviewById(ctx context.Context, id string, updatedReview models.NewReview) error {
	client := database.GetFirestoreClient()
	_, err := client.Collection(models.ReviewCollectionName).Doc(id).Update(ctx, []firestore.Update{
		{Path: "rating", Value: updatedReview.Rating},
		{Path: "content", Value: updatedReview.Content},
	})
	return err
}

func DeleteReviewById(ctx context.Context, id string) error {
	return common.DeleteItemById(ctx, models.ReviewCollectionName, id)
}
