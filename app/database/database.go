package database

import (
	"context"
	"fmt"
	"os"
	"scenic-spots-api/utils/logger"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var firestoreClient *firestore.Client

const SpotCollectionName string = "spots"
const ReviewCollectionName string = "reviews"
const UserAuthCollectionName string = "user_auth"

func InitializeFirestoreClient(ctx context.Context) error {
	var err error
	mode := os.Getenv("FIRESTORE_MODE")

	var connectFunc func(ctx context.Context) (*firestore.Client, error)
	if mode == "cloud" {
		connectFunc = connectToFirestoreCloud
	} else if mode == "emulator" {
		connectFunc = connectToEmulator
	} else {
		err = fmt.Errorf("invalid firestore mode - check config.go file")
		return err
	}

	firestoreClient, err = connectFunc(ctx)
	if err != nil {
		return err
	}

	if os.Getenv("DB_POPULATE") == "true" {
		if err = populateDatabase(ctx); err != nil {
			logger.Error(err.Error())
			return err
		}
	}

	logger.Success("Connected to firestore " + mode)
	return nil
}

func GetFirestoreClient() *firestore.Client {
	return firestoreClient
}

func connectToFirestoreCloud(ctx context.Context) (*firestore.Client, error) {
	pathToJson := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	projectId := os.Getenv("FIREBASE_PROJECT_ID")
	if pathToJson == "" {
		return nil, fmt.Errorf("FIREBASE_CREDENTIALS_PATH is not set - check .env file")
	}
	if projectId == "" {
		return nil, fmt.Errorf("FIREBASE_PROJECT_ID is not set - check .env file")
	}

	clientOption := option.WithCredentialsFile(pathToJson)
	client, err := firestore.NewClient(ctx, projectId, clientOption)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// USED FOR TESTING LOCALLY - MAKE SURE TO CONFIGURE THE EMULATOR ACCORDINGLY
func connectToEmulator(ctx context.Context) (*firestore.Client, error) {
	hostName := os.Getenv("FIRESTORE_EMULATOR_HOST_CONFIG")
	if hostName == "" {
		return nil, fmt.Errorf("FIRESTORE_EMULATOR_HOST_CONFIG is not set - check .env file")
	}
	os.Setenv("FIRESTORE_EMULATOR_HOST", hostName)
	projectID := os.Getenv("FIREBASE_PROJECT_ID")

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return client, nil
}
