package internal

import (
	"context"
	"net/http"
	"os"
	hHandler "scenic-spots-api/internal/api/handlers/health"
	sHandler "scenic-spots-api/internal/api/handlers/spot"
	uHandler "scenic-spots-api/internal/api/handlers/user"
	"scenic-spots-api/internal/database"
	"scenic-spots-api/utils/logger"

	"github.com/joho/godotenv"
)

func Start(ctx context.Context) error {
	if err := loadEnv(); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := database.InitializeFirestoreClient(ctx); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := database.InitalizeStorageClient(ctx); err != nil {
		logger.Error(err.Error())
		return err
	}
	initializeHandlers()
	return startTheServer()
}

func loadEnv() error {
	return godotenv.Load()
}

func initializeHandlers() {
	http.HandleFunc("/ping", hHandler.Ping)
	http.HandleFunc("/health", hHandler.Health)
	http.HandleFunc("/spot", sHandler.Spot)
	http.HandleFunc("/spot/", sHandler.SpotById)
	http.HandleFunc("/user/", uHandler.User)
}

func startTheServer() error {
	port := os.Getenv("PORT")
	logger.Info("Starting server on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}
