package database

import (
	"context"
	"fmt"
	"os"
	"scenic-spots-api/app/logger"
	"scenic-spots-api/configs"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var bucketHandle *storage.BucketHandle

func InitalizeStorageClient(ctx context.Context) error {
	var err error
	mode := configs.Env.StorageMode

	var connectFunc func(ctx context.Context) (*storage.BucketHandle, error)
	if mode == "cloud" {
		connectFunc = connectToStorageCloud
	} else if mode == "emulator" {
		connectFunc = connectToStorageEmulator
	} else {
		err = fmt.Errorf("invalid storage mode - check config.go file")
		logger.Error(err.Error())
		return err
	}

	bucketHandle, err = connectFunc(ctx)
	if err != nil {
		logger.Error("failed to connect to storage: " + err.Error())
		return err
	}

	logger.Success("Connected to storage " + mode)
	return nil
}

func connectToStorageCloud(ctx context.Context) (*storage.BucketHandle, error) {
	pathToJson := configs.Env.GoogleApplicationCredentials
	bucketName := configs.Env.StoragePhotoBucketName
	opt := option.WithCredentialsFile(pathToJson)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		logger.Error("error while initializing app")
		return nil, err
	}

	client, err := app.Storage(ctx)
	if err != nil {
		logger.Error("error getting Storage client")
		return nil, err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		logger.Error("error getting default bucket")
		return nil, err
	}

	return bucket, nil
}

func connectToStorageEmulator(ctx context.Context) (*storage.BucketHandle, error) {
	os.Setenv("STORAGE_EMULATOR_HOST", configs.Env.StorageEmulatorHost)
	bucketName := configs.Env.StoragePhotoBucketName

	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(bucketName)

	return bucket, nil
}

func GetStorageBucketHandle() *storage.BucketHandle {
	return bucketHandle
}
