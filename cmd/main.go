package main

import (
	"context"
	"os"
	"scenic-spots-api/app"
)

func main() {
	if err := app.Start(context.Background()); err != nil {
		os.Exit(1)
	}
}
