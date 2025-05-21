package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func RequestBodyIsEmpty(request *http.Request) bool {
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return false
	}
	if len(bodyBytes) > 0 {
		return false
	}

	return true
}

func DecodeAndValidateRequestBody[T any](request *http.Request, requestBodyStruct *T) error {
	if err := json.NewDecoder(request.Body).Decode(&requestBodyStruct); err != nil {
		return fmt.Errorf("Bad request body")
	}

	validate := validator.New()
	if err := validate.Struct(requestBodyStruct); err != nil {
		return fmt.Errorf("Invalid parameters")
	}
	return nil
}
