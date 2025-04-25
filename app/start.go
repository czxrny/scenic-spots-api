package app

import (
	"net/http"
	"os"
	"scenic-spots-api/app/database"
	"scenic-spots-api/app/handlers"
	"scenic-spots-api/app/logger"

	"github.com/joho/godotenv"
)

func Start() error {
	if err := loadEnv(); err != nil {
		logger.Error(err.Error())
		return err
	}
	database.InitializeFirestoreClient()
	initializeHandlers()
	return startTheServer()
}

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Couldn't load the .env file.")
		return err
	}
	return nil
}

func initializeHandlers() {
	http.HandleFunc("/ping", handlers.Ping)
	http.HandleFunc("/health", handlers.Health)
	http.HandleFunc("/spot", handlers.Spot)
	http.HandleFunc("/spot/", handlers.SpotById)
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
