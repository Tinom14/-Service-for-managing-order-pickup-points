package service_test

import (
	"avito_test/domain"
	"avito_test/repository"
	"avito_test/repository/mocks"
	"avito_test/usecases/service"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		role     string
		mockUser domain.User
		mockErr  error
		wantErr  bool
	}{
		{
			name:     "successful registration",
			email:    "test@example.com",
			password: "password123",
			role:     "user",
			mockUser: domain.User{Id: 1, Email: "test@example.com", Role: "user"},
			mockErr:  nil,
			wantErr:  false,
		},
		{
			name:     "registration error",
			email:    "test@example.com",
			password: "password123",
			role:     "user",
			mockUser: domain.User{},
			mockErr:  errors.New("repository error"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.User)
			mockRepo.On("Register", tt.email, mock.Anything, tt.role).Return(tt.mockUser, tt.mockErr)

			userService := service.NewUserService(mockRepo)
			user, err := userService.Register(tt.email, tt.password, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockUser, user)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	correctPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		email       string
		password    string
		mockUser    domain.User
		mockErr     error
		wantToken   bool
		wantErr     bool
		errContains string
	}{
		{
			name:      "successful login",
			email:     "test@example.com",
			password:  correctPassword,
			mockUser:  domain.User{Id: 1, Email: "test@example.com", Password: string(hashedPassword), Role: "user"},
			mockErr:   nil,
			wantToken: true,
			wantErr:   false,
		},
		{
			name:        "wrong email",
			email:       "wrong@example.com",
			password:    correctPassword,
			mockUser:    domain.User{},
			mockErr:     repository.NotFound,
			wantToken:   false,
			wantErr:     true,
			errContains: "wrong email",
		},
		{
			name:        "wrong password",
			email:       "test@example.com",
			password:    "wrongpassword",
			mockUser:    domain.User{Id: 1, Email: "test@example.com", Password: string(hashedPassword), Role: "user"},
			mockErr:     nil,
			wantToken:   false,
			wantErr:     true,
			errContains: "wrong password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.User)
			mockRepo.On("Login", tt.email).Return(tt.mockUser, tt.mockErr)

			userService := service.NewUserService(mockRepo)
			token, err := userService.Login(tt.email, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				if tt.wantToken {
					assert.NotEmpty(t, token)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
