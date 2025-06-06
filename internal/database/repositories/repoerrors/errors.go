package repoerrors

import (
	"errors"
)

// USED FOR /add METHOD FOR [/spot] ENDPOINT
var ErrAlreadyExists = errors.New("item already exists")

// USED FOR /get METHOD FOR [/spot/{id} & /spot/{id}/review/{rId}] ENDPOINTS
var ErrDoesNotExist = errors.New("item does not exist")

// ==============================================
