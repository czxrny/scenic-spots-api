package common

import (
	"context"
	"scenic-spots-api/app/database"
	"scenic-spots-api/internal/repoerrors"
	"scenic-spots-api/models"
	"scenic-spots-api/utils/generics"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func FindItemById[T models.Identifiable](ctx context.Context, collectionName string, id string) (T, error) {
	client := database.GetFirestoreClient()
	docRef := client.Collection(collectionName).Doc(id)
	doc, err := docRef.Get(ctx)
	if err != nil {
		var zero T
		if !doc.Exists() {
			return zero, repoerrors.ErrDoesNotExist
		}
		return zero, err
	}

	var item T
	if err := doc.DataTo(&item); err != nil {
		var zero T
		return zero, err
	}

	item.SetId(docRef.ID)
	return item, nil
}

func GetAllItems[T models.Identifiable](ctx context.Context, query firestore.Query) ([]T, error) {
	iter := query.Documents(ctx)
	found := make([]T, 0)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []T{}, err
		}

		var item T

		if err := doc.DataTo(&item); err != nil {
			return []T{}, err
		}

		item.SetId(doc.Ref.ID)

		found = append(found, item)
	}
	return found, nil
}

func AddItem[T models.Identifiable](ctx context.Context, collectionName string, item T) (T, error) {
	// Casting to a json to avoid capitalized words in database.
	data, err := generics.StructToMapLower(item)
	if err != nil {
		var zero T
		return zero, err
	}

	client := database.GetFirestoreClient()
	collectionRef := client.Collection(collectionName)
	docRef, _, err := collectionRef.Add(ctx, data)

	if err != nil {
		var zero T
		return zero, err
	}

	item.SetId(docRef.ID)

	return item, nil
}

func DeleteItemById(ctx context.Context, collectionName string, id string) error {
	client := database.GetFirestoreClient()
	docRef := client.Collection(collectionName).Doc(id)

	if _, err := docRef.Delete(ctx); err != nil {
		return err
	}

	return nil
}

func DeleteAllItems(ctx context.Context, query firestore.Query) error {
	iter := query.Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		doc.Ref.Delete(ctx)
	}
	return nil
}
