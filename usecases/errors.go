package usecases

import "errors"

var (
	ErrUnclosedReception = errors.New("unclosed reception")
	ErrAlreadyClosed     = errors.New("already closed")
)
