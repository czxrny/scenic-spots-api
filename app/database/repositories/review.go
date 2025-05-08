package repositories

import (
	"context"
	"scenic-spots-api/app/database"
	"scenic-spots-api/models"
	"strconv"
	"time"

	"google.golang.org/api/iterator"
)

var reviewCollectionName string = "reviews"

func GetReviews(ctx context.Context, spotId string, limitParam string) ([]models.Review, error) {
	if _, err := FindById(ctx, spotId); err != nil {
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
	if _, err := FindById(ctx, spotId); err != nil {
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
