package repositories

import (
	"context"
	"fmt"
	"scenic-spots-api/app/database"
	"scenic-spots-api/models"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const reviewCollectionName string = "reviews"

func GetReviews(ctx context.Context, spotId string, limitParam string) ([]models.Review, error) {
	if _, err := FindSpotById(ctx, spotId); err != nil {
		return []models.Review{}, err
	}

	var limit int
	var err error

	client := database.GetFirestoreClient()
	query := client.Collection(reviewCollectionName).Query.Where("spotId", "==", spotId)

	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return []models.Review{}, err
		}
		query = query.Limit(limit)
	}

	iter := query.Documents(ctx)

	found := make([]models.Review, 0)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []models.Review{}, err
		}

		var review models.Review

		if err := doc.DataTo(&review); err != nil {
			return []models.Review{}, err
		}

		found = append(found, models.Review{
			Id:        doc.Ref.ID,
			SpotId:    spotId,
			Rating:    review.Rating,
			Content:   review.Content,
			AddedBy:   review.AddedBy,
			CreatedAt: review.CreatedAt,
		})
	}
	return found, nil
}

func AddReview(ctx context.Context, spotId string, reviewInfo models.NewReview) ([]models.Review, error) {
	if _, err := FindSpotById(ctx, spotId); err != nil {
		return []models.Review{}, err
	}

	client := database.GetFirestoreClient()
	collectionRef := client.Collection(reviewCollectionName)

	review := models.Review{
		SpotId:    spotId,
		Rating:    reviewInfo.Rating,
		Content:   reviewInfo.Content,
		AddedBy:   "test user", /* TODO */
		CreatedAt: time.Now(),
	}

	// Casting to a json to avoid capitalized words in database.
	reviewData := map[string]interface{}{
		"spotId":      review.SpotId,
		"description": review.Rating,
		"content":     review.Content,
		"addedBy":     review.AddedBy,
		"createdAt":   review.CreatedAt,
	}

	docRef, _, err := collectionRef.Add(ctx, reviewData)

	if err != nil {
		return []models.Review{}, err
	}

	review.Id = docRef.ID
	var result []models.Review

	result = append(result, review)

	return result, nil
}

func FindReviewById(ctx context.Context, id string) ([]models.Review, error) {
	client := database.GetFirestoreClient()
	doc, err := client.Collection(reviewCollectionName).Doc(id).Get(ctx)
	if err != nil {
		return []models.Review{}, err
	}

	var review models.Review
	if err := doc.DataTo(&review); err != nil {
		return []models.Review{}, err
	}

	// Content won't ever be empty - unless no doc was found.
	if review.Content == "" {
		return []models.Review{}, fmt.Errorf("Spot with ID: %v does not exist", id)
	}

	review.Id = doc.Ref.ID
	var result []models.Review

	result = append(result, review)

	return result, nil
}

func UpdateReviewById(ctx context.Context, id string, newValues models.NewReview) ([]models.Review, error) {
	var err error
	result := []models.Review{}

	result, err = FindReviewById(ctx, id)
	if err != nil {
		return []models.Review{}, err
	}

	reviewToUpdate := result[0]
	if newValues.Rating != reviewToUpdate.Rating {
		reviewToUpdate.Rating = newValues.Rating
	}
	if newValues.Content != "" {
		reviewToUpdate.Content = newValues.Content
	}

	client := database.GetFirestoreClient()
	_, err = client.Collection(reviewCollectionName).Doc(id).Update(ctx, []firestore.Update{
		{Path: "rating", Value: reviewToUpdate.Rating},
		{Path: "content", Value: reviewToUpdate.Content},
	})
	if err != nil {
		return []models.Review{}, err
	}

	result[0] = reviewToUpdate

	return result, nil
}

func DeleteReviewById(ctx context.Context, id string) error {
	client := database.GetFirestoreClient()
	docRef := client.Collection(reviewCollectionName).Doc(id)

	_, err := docRef.Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
