package database

import (
	"context"
	"fmt"
	"os"
	"scenic-spots-api/app/logger"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var firestoreClient *firestore.Client

func InitializeFirestoreClient() error {
	var err error
	ctx := context.Background()
	mode := os.Getenv("FIRESTORE_MODE")

	var connectFunc func(ctx context.Context) (*firestore.Client, error)
	if mode == "cloud" {
		connectFunc = connectToFirestoreCloud
	} else if mode == "emulator" {
		connectFunc = connectToEmulator
	} else {
		err = fmt.Errorf("invalid firestore mode - check .env file")
		logger.Error(err.Error())
		return err
	}

	firestoreClient, err = connectFunc(ctx)
	if err != nil {
		logger.Error("failed to connect to firestore: " + err.Error())
		return err
	}

	logger.Success("Connected to firestore " + mode)
	return nil
}

func GetFirestoreClient() *firestore.Client {
	return firestoreClient
}

func connectToFirestoreCloud(ctx context.Context) (*firestore.Client, error) {
	pathToJson := os.Getenv("FIRESTORE_CREDENTIALS_PATH")
	projectId := os.Getenv("FIRESTORE_PROJECT_ID")
	clientOption := option.WithCredentialsFile(pathToJson)

	client, err := firestore.NewClient(ctx, projectId, clientOption)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// USED FOR TESTING LOCALLY - MAKE SURE TO CONFIGURE THE EMULATOR ACCORDINGLY
func connectToEmulator(ctx context.Context) (*firestore.Client, error) {
	projectID := "demo"

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return client, nil
}
