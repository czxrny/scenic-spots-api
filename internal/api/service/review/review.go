package review

import (
	"context"
	"net/url"
	"scenic-spots-api/internal/auth"
	reviewRepo "scenic-spots-api/internal/database/repositories/review"
	spotRepo "scenic-spots-api/internal/database/repositories/spot"
	"scenic-spots-api/internal/models"
	"time"
)

func GetReview(ctx context.Context, spotId string, query url.Values) ([]models.Review, error) {
	if _, err := spotRepo.FindSpotById(ctx, spotId); err != nil {
		return []models.Review{}, err
	}

	params := models.ReviewQueryParams{
		SpotId: spotId,
		Limit:  query.Get("limit"),
	}

	found, err := reviewRepo.GetReviews(ctx, params)
	if err != nil {
		return []models.Review{}, err
	}
	return found, nil
}

func AddReview(ctx context.Context, token string, newReviewInfo models.NewReview) ([]models.Review, error) {
	// Check if the spot exists!
	if _, err := spotRepo.FindSpotById(ctx, newReviewInfo.SpotId); err != nil {
		return []models.Review{}, err
	}

	userName, err := auth.ExtractFromToken(token, "usr")
	if err != nil {
		return []models.Review{}, err
	}

	review := models.Review{
		SpotId:    newReviewInfo.SpotId,
		Rating:    newReviewInfo.Rating,
		Content:   newReviewInfo.Content,
		AddedBy:   userName,
		CreatedAt: time.Now(),
	}

	addedReview, err := reviewRepo.AddReview(ctx, review)
	if err != nil {
		return []models.Review{}, err
	}

	var result []models.Review
	result = append(result, addedReview)

	return result, nil
}

func FindReviewById(ctx context.Context, id string) ([]models.Review, error) {
	review, err := reviewRepo.FindReviewById(ctx, id)
	if err != nil {
		return []models.Review{}, err
	}

	var result []models.Review
	result = append(result, review)

	return result, nil
}

func UpdateReviewById(ctx context.Context, token string, newReviewInfo models.NewReview, reviewId string) ([]models.Review, error) {
	review, err := reviewRepo.FindReviewById(ctx, reviewId)
	if err != nil {
		return []models.Review{}, err
	}

	if err := auth.IsAuthorizedToEditAsset(token, review.AddedBy); err != nil {
		return []models.Review{}, err
	}

	if err := reviewRepo.UpdateReviewById(ctx, reviewId, newReviewInfo); err != nil {
		return []models.Review{}, err
	}

	review.Rating = newReviewInfo.Rating
	review.Content = newReviewInfo.Content

	var result []models.Review
	result = append(result, review)

	return result, nil
}

func DeleteReviewById(ctx context.Context, token string, reviewId string) error {
	review, err := reviewRepo.FindReviewById(ctx, reviewId)
	if err != nil {
		return err
	}

	if err := auth.IsAuthorizedToEditAsset(token, review.AddedBy); err != nil {
		return err
	}

	return reviewRepo.DeleteReviewById(ctx, reviewId)
}

func DeleteAllReviews(ctx context.Context, token string, spotId string) error {
	// can delete only if jwt states that the user is an admin.
	if err := auth.IsAuthorizedToEditAsset(token, ""); err != nil {
		return err
	}

	return reviewRepo.DeleteReviewById(ctx, spotId)
}
