package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

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
