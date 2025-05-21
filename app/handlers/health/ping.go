package health

import (
	"net/http"
	"scenic-spots-api/utils/logger"
)

func Ping(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	logger.Info("Ping request.")
}
