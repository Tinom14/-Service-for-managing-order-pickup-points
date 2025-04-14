package types

import (
	"avito_test/repository"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrInvalidJSON           = errors.New("invalid json")
	ErrEmailPasswordRequired = errors.New("email and password are required")
	ErrInvalidEmail          = errors.New("invalid email")
	ErrInvalidRole           = errors.New("invalid role")
	ErrInvalidCity           = errors.New("invalid city")
	ErrPvzIdRequired         = errors.New("pvzId is required")
	ErrTypePvzIdRequired     = errors.New("type and pvzId are required")
)

func AuthError(w http.ResponseWriter, err error, resp any) {
	if errors.Is(err, repository.NotFound) {
		http.Error(w, "invalid email", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}
