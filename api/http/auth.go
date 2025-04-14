package http

import (
	"avito_test/api/http/types"
	"avito_test/repository"
	"avito_test/usecases"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type User struct {
	Service usecases.User
}

func NewUserHandler(service usecases.User) *User {
	return &User{Service: service}
}

func (u *User) DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateDummyLoginHandlerRequest(r)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrInvalidJSON):
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		case errors.Is(err, types.ErrInvalidRole):
			http.Error(w, "Invalid role", http.StatusBadRequest)
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		return
	}

	token, err := u.Service.GetToken("1", req.Role)

	types.AuthError(w, err, types.LoginHandlerResponse{Token: token})
}

func (u *User) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateRegisterHandlerRequest(r)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrInvalidEmail):
			http.Error(w, "Invalid email format", http.StatusBadRequest)
		case errors.Is(err, types.ErrEmailPasswordRequired):
			http.Error(w, "Email and password are required", http.StatusBadRequest)
		case errors.Is(err, types.ErrInvalidJSON):
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		case errors.Is(err, types.ErrInvalidRole):
			http.Error(w, "Invalid role", http.StatusBadRequest)
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		return
	}

	user, err := u.Service.Register(req.Email, req.Password, req.Role)
	if errors.Is(err, repository.ErrEmailAlreadyExists) {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}

func (u *User) LoginHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateLoginHandlerRequest(r)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := u.Service.Login(req.Email, req.Password)

	types.AuthError(w, err, types.LoginHandlerResponse{Token: token})
}

func (u *User) WithUserHandlers(r chi.Router) {
	r.Post("/dummyLogin", u.DummyLoginHandler)
	r.Post("/register", u.RegisterHandler)
	r.Post("/login", u.LoginHandler)
}
