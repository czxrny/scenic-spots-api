package repositories

import (
	"context"
	"fmt"
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
		return zero, err
	}

	if !doc.Exists() {
		var zero T
		return zero, fmt.Errorf("document with ID: [%v] in the collection [%v] does not exist", id, collectionName)
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

//func addItem

// func updateById[T models.Identifiable, T1 any](ctx context.Context, collectionName string, id string,) (T, error) {
// 	itemToUpdate, err := findById[T](ctx, collectionName, id)
// 	if err != nil {
// 		var zero T
// 		return zero, err
// 	}
// }

func deleteItemById(ctx context.Context, collectionName string, id string) error {
	client := database.GetFirestoreClient()
	docRef := client.Collection(collectionName).Doc(id)

	_, err := docRef.Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
