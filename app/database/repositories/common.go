package repositories

import (
	"context"
	"scenic-spots-api/app/database"
	"scenic-spots-api/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func findItemById[T models.Identifiable](ctx context.Context, collectionName string, id string) (T, error) {
	client := database.GetFirestoreClient()
	docRef := client.Collection(collectionName).Doc(id)
	doc, err := docRef.Get(ctx)
	if err != nil {
		var zero T
		if !doc.Exists() {
			return zero, ErrDoesNotExist
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

func getAllItems[T models.Identifiable](ctx context.Context, query firestore.Query) ([]T, error) {
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

func deleteItemById(ctx context.Context, collectionName string, id string) error {
	client := database.GetFirestoreClient()
	docRef := client.Collection(collectionName).Doc(id)

	if _, err := docRef.Delete(ctx); err != nil {
		return err
	}

	return nil
}

func deleteAllItems(ctx context.Context, query firestore.Query) error {
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
