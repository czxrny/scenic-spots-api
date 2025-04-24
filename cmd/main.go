package main

import (
	"os"
	"scenic-spots-api/app"
)

func main() {
	var port string = "8080"

	if err := app.Start(port); err != nil {
		os.Exit(1)
	}
}
