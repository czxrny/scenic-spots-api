package apierrors

import (
	"errors"
	"fmt"
)

var ErrInvalidCredentials = errors.New("invalid username or password")
var ErrIsUnauthorized = errors.New("user is unauthorized to edit the asset")

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
