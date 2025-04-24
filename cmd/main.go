package main

import (
	"scenic-spots-api/app"
)

func main() {
	var port string = "8080"

	app.Start(port)
}
