package api

import "errors"

var ErrIsUnauthorized = errors.New("user is unauthorized to edit the asset")
