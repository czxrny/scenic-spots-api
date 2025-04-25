package main

import (
	"os"
	"scenic-spots-api/app"
)

func main() {
	if err := app.Start(); err != nil {
		os.Exit(1)
	}
}
