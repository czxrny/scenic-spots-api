package app

import (
	"context"
	"net/http"
	"scenic-spots-api/app/database"
	"scenic-spots-api/app/handlers"
	"scenic-spots-api/app/logger"
	"scenic-spots-api/configs"
)

func Start(ctx context.Context) error {
	configs.InitializeVariables()
	if err := database.InitializeFirestoreClient(ctx); err != nil {
		return err
	}
	initializeHandlers()
	return startTheServer()
}

func initializeHandlers() {
	http.HandleFunc("/ping", handlers.Ping)
	http.HandleFunc("/health", handlers.Health)
	http.HandleFunc("/spot", handlers.Spot)
	http.HandleFunc("/spot/", handlers.SpotById)
}

func startTheServer() error {
	port := configs.Env.Port
	logger.Info("Starting server on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}
