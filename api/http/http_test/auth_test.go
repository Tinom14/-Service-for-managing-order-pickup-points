package http_test

import (
	http2 "avito_test/api/http"
	"avito_test/domain"
	"avito_test/repository"
	"avito_test/usecases/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_DummyLogin(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		mockSetup    func(*mocks.User)
		expectedCode int
	}{
		{
			name:        "Success employee login",
			requestBody: `{"role": "employee"}`,
			mockSetup: func(m *mocks.User) {
				m.On("GetToken", "1", "employee").Return("employee-token", nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid role",
			requestBody:  `{"role": "invalid"}`,
			mockSetup:    func(m *mocks.User) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.User)
			tt.mockSetup(mockService)
			handler := http2.NewUserHandler(mockService)

			req := httptest.NewRequest("POST", "/dummyLogin", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/dummyLogin", handler.DummyLoginHandler)
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedCode == http.StatusOK {
				var response struct{ Token string }
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "employee-token", response.Token)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_DummyLogin_InvalidJSON(t *testing.T) {
	handler := &http2.User{Service: new(mocks.User)}

	req := httptest.NewRequest("POST", "/dummyLogin", bytes.NewBufferString("{invalid"))
	rec := httptest.NewRecorder()

	handler.DummyLoginHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_Register(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		mockSetup    func(*mocks.User)
		expectedCode int
	}{
		{
			name:        "Success registration",
			requestBody: `{"email": "test@test.com", "password": "password", "role": "employee"}`,
			mockSetup: func(m *mocks.User) {
				m.On("Register", "test@test.com", mock.Anything, "employee").
					Return(domain.User{Id: 1, Email: "test@test.com", Role: "employee"}, nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:        "Duplicate email",
			requestBody: `{"email": "exists@test.com", "password": "password", "role": "employee"}`,
			mockSetup: func(m *mocks.User) {
				m.On("Register", "exists@test.com", mock.Anything, "employee").
					Return(domain.User{}, repository.ErrEmailAlreadyExists)
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.User)
			tt.mockSetup(mockService)
			handler := http2.NewUserHandler(mockService)

			req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/register", handler.RegisterHandler)
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		mockSetup    func(*mocks.User)
		expectedCode int
	}{
		{
			name:        "Success login",
			requestBody: `{"email": "test@test.com", "password": "password"}`,
			mockSetup: func(m *mocks.User) {
				m.On("Login", "test@test.com", "password").Return("valid-token", nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Invalid credentials",
			requestBody: `{"email": "test@test.com", "password": "wrong"}`,
			mockSetup: func(m *mocks.User) {
				m.On("Login", "test@test.com", "wrong").Return("", errors.New("invalid credentials"))
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.User)
			tt.mockSetup(mockService)
			handler := http2.NewUserHandler(mockService)

			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/login", handler.LoginHandler)
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}
