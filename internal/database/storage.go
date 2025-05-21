package database

import (
	"context"
	"fmt"
	"os"
	"scenic-spots-api/utils/logger"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var bucketHandle *storage.BucketHandle
var bucketName string

func InitalizeStorageClient(ctx context.Context) error {
	var err error
	mode := os.Getenv("STORAGE_MODE")

	var connectFunc func(ctx context.Context) (*storage.BucketHandle, error)
	if mode == "cloud" {
		connectFunc = connectToStorageCloud
	} else if mode == "emulator" {
		connectFunc = connectToStorageEmulator
	} else {
		err = fmt.Errorf("invalid storage mode - check .env file")
		return err
	}

	setBucketName()
	bucketHandle, err = connectFunc(ctx)
	if err != nil {
		return err
	}

	logger.Success("Connected to storage " + mode)
	return nil
}

func connectToStorageCloud(ctx context.Context) (*storage.BucketHandle, error) {
	pathToJson := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	opt := option.WithCredentialsFile(pathToJson)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

func connectToStorageEmulator(ctx context.Context) (*storage.BucketHandle, error) {
	hostName := os.Getenv("STORAGE_EMULATOR_HOST_CONFIG")
	if hostName == "" {
		return nil, fmt.Errorf("STORAGE_EMULATOR_HOST_CONFIG is not set - check .env file")
	}
	os.Setenv("STORAGE_EMULATOR_HOST", hostName)

	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}

	return client.Bucket(bucketName), nil
}

func setBucketName() {
	bucketName = os.Getenv("STORAGE_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "default"
	}
}

func GetStorageBucketHandle() *storage.BucketHandle {
	return bucketHandle
}
