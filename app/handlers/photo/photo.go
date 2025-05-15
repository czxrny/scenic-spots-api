package photo

import (
	"net/http"
	"scenic-spots-api/app/logger"
)

func getPhoto(response http.ResponseWriter, request *http.Request, id string) {
	logger.Info("Get photo for the spot")
}

func addPhoto(response http.ResponseWriter, request *http.Request, id string) {
	logger.Info("Add photo to the spot")
}

func deletePhoto(response http.ResponseWriter, request *http.Request, spotId string, photoId string) {
	logger.Info("Delete specified photo of the spot")
}

func Photo(response http.ResponseWriter, request *http.Request, id string) {
	method := request.Method

	switch method {
	case "GET":
		getPhoto(response, request, id)
	case "POST":
		addPhoto(response, request, id)
	case "DELETE":
		// check if there is photoId specified! if not, delete all!
		deletePhoto(response, request, id, "22")
	}
}
