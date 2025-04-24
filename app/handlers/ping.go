package handlers

import (
	"fmt"
	"net/http"
)

func Ping(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	fmt.Println("Ping request.")
}
