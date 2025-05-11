package repositories

import (
	"errors"
	"fmt"
)

// USED FOR /add METHOD FOR [/spot] ENDPOINT
var ErrSpotAlreadyExists = errors.New("item already exists")

// USED FOR /get METHOD FOR [/spot/{id} & /spot/{id}/review/{rId}] ENDPOINTS
var ErrDoesNotExist = errors.New("item does not exist")

// USED FOR /get METHODS WITH QUERY PARAMS - ALL INVALID PARAMETER ERRORS FALL INTO ErrInvalidSpotParameters
var ErrInvalidQueryParameters = fmt.Errorf("invalid query parameters")

type InvalidQueryParameterError struct {
	Message string
}

// Error message
func (e *InvalidQueryParameterError) Error() string {
	return e.Message
}

// For errors.Is func: InvalidSpotParameterError is wrapped with ErrInvalidSpotParameters error
func (e *InvalidQueryParameterError) Unwrap() error {
	return ErrInvalidQueryParameters
}

// ==============================================
