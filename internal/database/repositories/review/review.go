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
	"time"

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
	if _, err := common.FindItemById[*models.Spot](ctx, database.SpotCollectionName, params.SpotId); err != nil {
		return []models.Review{}, err
	}

	client := database.GetFirestoreClient()
	collectionRef := client.Collection(database.ReviewCollectionName)

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

func AddReview(ctx context.Context, reviewInfo models.NewReview) ([]models.Review, error) {
	if _, err := common.FindItemById[*models.Spot](ctx, database.SpotCollectionName, reviewInfo.SpotId); err != nil {
		return []models.Review{}, err
	}

	review := models.Review{
		SpotId:    reviewInfo.SpotId,
		Rating:    reviewInfo.Rating,
		Content:   reviewInfo.Content,
		AddedBy:   "test user", /* TODO */
		CreatedAt: time.Now(),
	}

	addedReview, err := common.AddItem(ctx, database.SpotCollectionName, &review)
	if err != nil {
		return []models.Review{}, err
	}

	var result []models.Review
	result = append(result, *addedReview)

	return result, nil
}

func DeleteAllReviews(ctx context.Context, spotId string) error {
	client := database.GetFirestoreClient()
	query := client.Collection(database.ReviewCollectionName).Where("spotId", "==", spotId)

	return common.DeleteAllItems(ctx, query)
}

func FindReviewById(ctx context.Context, id string) ([]models.Review, error) {
	review, err := common.FindItemById[*models.Review](ctx, database.ReviewCollectionName, id)
	if err != nil {
		return []models.Review{}, err
	}

	var result []models.Review
	result = append(result, *review)

	return result, nil
}

func UpdateReviewById(ctx context.Context, id string, updatedReview models.Review) ([]models.Review, error) {
	var err error
	result := []models.Review{}

	client := database.GetFirestoreClient()
	_, err = client.Collection(database.ReviewCollectionName).Doc(id).Update(ctx, []firestore.Update{
		{Path: "rating", Value: updatedReview.Rating},
		{Path: "content", Value: updatedReview.Content},
	})
	if err != nil {
		return []models.Review{}, err
	}

	result = append(result, updatedReview)

	return result, nil
}

func DeleteReviewById(ctx context.Context, id string) error {
	return common.DeleteItemById(ctx, database.ReviewCollectionName, id)
}
