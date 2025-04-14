package repository

import "errors"

var (
	NotFound              = errors.New("not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
