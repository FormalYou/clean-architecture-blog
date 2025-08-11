package repository

import "errors"

// ErrNotFound is a common error for when a record is not found in the repository.
var ErrNotFound = errors.New("not found")
