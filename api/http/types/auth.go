package types

import (
	"avito_test/domain"
	"encoding/json"
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"net/http"
)

type DummyLoginHandlerRequest struct {
	Role string `json:"role"`
}

func CreateDummyLoginHandlerRequest(r *http.Request) (*DummyLoginHandlerRequest, error) {
	var req DummyLoginHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, ErrInvalidJSON
	}
	if req.Role != "moderator" && req.Role != "employee" {
		return nil, ErrInvalidRole
	}
	return &req, nil
}

type LoginHandlerResponse struct {
	Token string `json:"token"`
}

type RegisterHandlerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func CreateRegisterHandlerRequest(r *http.Request) (*RegisterHandlerRequest, error) {
	var req RegisterHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, ErrInvalidJSON
	}
	if req.Email == "" || req.Password == "" {
		return nil, ErrEmailPasswordRequired
	}

	if req.Role != "moderator" && req.Role != "employee" {
		return nil, ErrInvalidRole
	}

	if !ev.NewSyntaxValidator().Validate(ev.NewInput(evmail.FromString(req.Email))).IsValid() {
		return nil, ErrInvalidEmail
	}
	return &req, nil
}

type RegisterHandlerResponse struct {
	User domain.User
}

type LoginHandlerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateLoginHandlerRequest(r *http.Request) (*LoginHandlerRequest, error) {
	var req LoginHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, ErrInvalidJSON
	}
	if req.Email == "" || req.Password == "" {
		return nil, ErrEmailPasswordRequired
	}
	return &req, nil
}
