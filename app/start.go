package app

import (
	"net/http"
	"scenic-spots-api/app/handlers"
	"scenic-spots-api/app/logger"
)

func Start(port string) error {
	initializeHandlers()
	return startTheServer(port)
}

func initializeHandlers() {
	http.HandleFunc("/ping", handlers.Ping)
	http.HandleFunc("/health", handlers.Health)
}

func startTheServer(port string) error {
	logger.Info("Starting server on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}
