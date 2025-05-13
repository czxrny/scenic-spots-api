package repositories

import (
	"context"
	"fmt"
	"scenic-spots-api/app/database"
	"scenic-spots-api/models"
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
	if _, err := findItemById[*models.Spot](ctx, database.SpotCollectionName, params.SpotId); err != nil {
		return []models.Review{}, err
	}

	client := database.GetFirestoreClient()
	collectionRef := client.Collection(database.ReviewCollectionName)

	query, err := buildReviewQuery(collectionRef, params)
	if err != nil {
		return []models.Review{}, &InvalidQueryParameterError{
			Message: err.Error(),
		}
	}

	found, err := getAllItems[*models.Review](ctx, query)
	if err != nil {
		return []models.Review{}, err
	}

	result := generics.DereferenceAll(found)
	return result, nil
}

func AddReview(ctx context.Context, reviewInfo models.NewReview) ([]models.Review, error) {
	if _, err := findItemById[*models.Spot](ctx, database.SpotCollectionName, reviewInfo.SpotId); err != nil {
		return []models.Review{}, err
	}

	review := models.Review{
		SpotId:    reviewInfo.SpotId,
		Rating:    reviewInfo.Rating,
		Content:   reviewInfo.Content,
		AddedBy:   "test user", /* TODO */
		CreatedAt: time.Now(),
	}

	// Casting to a json to avoid capitalized words in database.
	reviewData := map[string]interface{}{
		"spotId":    review.SpotId,
		"rating":    review.Rating,
		"content":   review.Content,
		"addedBy":   review.AddedBy,
		"createdAt": review.CreatedAt,
	}

	client := database.GetFirestoreClient()
	collectionRef := client.Collection(database.ReviewCollectionName)
	docRef, _, err := collectionRef.Add(ctx, reviewData)

	if err != nil {
		return []models.Review{}, err
	}

	review.Id = docRef.ID
	var result []models.Review

	result = append(result, review)

	return result, nil
}

func DeleteAllReviews(ctx context.Context, spotId string) error {
	client := database.GetFirestoreClient()
	query := client.Collection(database.ReviewCollectionName).Where("spotId", "==", spotId)

	return deleteAllItems(ctx, query)
}

func FindReviewById(ctx context.Context, id string) ([]models.Review, error) {
	review, err := findItemById[*models.Review](ctx, database.ReviewCollectionName, id)
	if err != nil {
		return []models.Review{}, err
	}

	var result []models.Review
	result = append(result, *review)

	return result, nil
}

func UpdateReviewById(ctx context.Context, id string, newValues models.NewReview) ([]models.Review, error) {
	var err error
	result := []models.Review{}

	reviewToUpdate, err := findItemById[*models.Review](ctx, database.ReviewCollectionName, id)
	if err != nil {
		return []models.Review{}, err
	}

	if newValues.Rating != reviewToUpdate.Rating {
		reviewToUpdate.Rating = newValues.Rating
	}
	if newValues.Content != "" {
		reviewToUpdate.Content = newValues.Content
	}

	client := database.GetFirestoreClient()
	_, err = client.Collection(database.ReviewCollectionName).Doc(id).Update(ctx, []firestore.Update{
		{Path: "rating", Value: reviewToUpdate.Rating},
		{Path: "content", Value: reviewToUpdate.Content},
	})
	if err != nil {
		return []models.Review{}, err
	}

	result = append(result, *reviewToUpdate)

	return result, nil
}

func DeleteReviewById(ctx context.Context, id string) error {
	return deleteItemById(ctx, database.ReviewCollectionName, id)
}
