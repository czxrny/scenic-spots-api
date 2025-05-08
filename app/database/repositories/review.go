package repositories

import (
	"context"
	"scenic-spots-api/app/database"
	"scenic-spots-api/models"
	"strconv"

	"google.golang.org/api/iterator"
)

var reviewCollectionName string = "reviews"

func GetReviews(ctx context.Context, spotId string, limitParam string) ([]models.Review, error) {
	if _, err := FindById(ctx, spotId); err != nil {
		return []models.Review{}, err
	}

	var limit int
	var err error

	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return []models.Review{}, err
		}
	}

	client := database.GetFirestoreClient()
	iter := client.Collection(reviewCollectionName).Query.Limit(limit).Documents(ctx)

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
