package health

import (
	"fmt"
	"math/rand"
	"net/http"
	"scenic-spots-api/utils/logger"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

func Health(response http.ResponseWriter, request *http.Request) {
	var meals [2]string = [2]string{"pizza", "ramen"}
	var uptime time.Duration = time.Since(startTime)
	var decision int = rand.Intn(2)

	response.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(response, `{"status": "ok, is alive!", "uptime": "%s", "would really like some": "%s :)"}`, uptime, meals[decision])
	logger.Info("Health request")
}
