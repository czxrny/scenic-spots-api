package database

import (
	"context"
	"fmt"
	"scenic-spots-api/app/logger"
	"scenic-spots-api/configs"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var firestoreClient *firestore.Client

func InitializeFirestoreClient(ctx context.Context) error {
	var err error
	mode := configs.Env.FirestoreMode

	var connectFunc func(ctx context.Context) (*firestore.Client, error)
	if mode == "cloud" {
		connectFunc = connectToFirestoreCloud
	} else if mode == "emulator" {
		connectFunc = connectToEmulator
	} else {
		err = fmt.Errorf("invalid firestore mode - check config.go file")
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
	pathToJson := configs.Env.GoogleApplicationCredentials
	projectId := configs.Env.FirestoreProjectID
	if pathToJson == "" {
		return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS is not set - check config.go file")
	}
	if projectId == "" {
		return nil, fmt.Errorf("FIRESTORE_PROJECT_ID is not set - check config.go file")
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
	projectID := "demo"

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return client, nil
}
