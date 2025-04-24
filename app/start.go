package app

import (
	"fmt"
	"net/http"
	"scenic-spots-api/app/handlers"
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
	fmt.Printf("Starting server on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Server started successfully")
	return nil
}
