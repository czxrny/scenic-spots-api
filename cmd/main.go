package main

import (
	"context"
	"os"
	app "scenic-spots-api/internal"
)

func main() {
	if err := app.Start(context.Background()); err != nil {
		os.Exit(1)
	}
}
