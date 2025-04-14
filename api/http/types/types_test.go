package types

import (
	"avito_test/repository"
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateDummyLoginHandlerRequest(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantErr  bool
		expected DummyLoginHandlerRequest
	}{
		{
			name:    "Valid employee role",
			body:    `{"role": "employee"}`,
			wantErr: false,
			expected: DummyLoginHandlerRequest{
				Role: "employee",
			},
		},
		{
			name:    "Invalid role",
			body:    `{"role": "invalid"}`,
			wantErr: true,
		},
		{
			name:    "Empty body",
			body:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/dummyLogin", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			got, err := CreateDummyLoginHandlerRequest(req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, *got)
		})
	}
}

func TestCreateRegisterHandlerRequest(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantErr  bool
		expected RegisterHandlerRequest
	}{
		{
			name:    "Valid request",
			body:    `{"email": "test@test.com", "password": "pass", "role": "employee"}`,
			wantErr: false,
			expected: RegisterHandlerRequest{
				Email:    "test@test.com",
				Password: "pass",
				Role:     "employee",
			},
		},
		{
			name:    "Invalid email",
			body:    `{"email": "invalid", "password": "pass", "role": "employee"}`,
			wantErr: true,
		},
		{
			name:    "Missing password",
			body:    `{"email": "test@test.com", "role": "employee"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			got, err := CreateRegisterHandlerRequest(req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, *got)
		})
	}
}

func TestCreateLoginHandlerRequest(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantErr  bool
		expected LoginHandlerRequest
	}{
		{
			name:    "Valid request",
			body:    `{"email": "test@test.com", "password": "pass"}`,
			wantErr: false,
			expected: LoginHandlerRequest{
				Email:    "test@test.com",
				Password: "pass",
			},
		},
		{
			name:    "Empty email",
			body:    `{"password": "pass"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			got, err := CreateLoginHandlerRequest(req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, *got)
		})
	}
}

func TestCreateOpenPvzHandlerRequest(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantErr  bool
		expected OpenPvzHandlerRequest
	}{
		{
			name:    "Valid city",
			body:    `{"city": "Москва"}`,
			wantErr: false,
			expected: OpenPvzHandlerRequest{
				City: "Москва",
			},
		},
		{
			name:    "Invalid city",
			body:    `{"city": "Новосибирск"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/pvz", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			got, err := CreateOpenPvzHandlerRequest(req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, *got)
		})
	}
}

func TestCreateStartReceptionHandlerRequest(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantErr  bool
		expected StartReceptionHandlerRequest
	}{
		{
			name:    "Valid request",
			body:    `{"pvzId": "123"}`,
			wantErr: false,
			expected: StartReceptionHandlerRequest{
				PvzId: "123",
			},
		},
		{
			name:    "Missing pvzId",
			body:    `{}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/receptions", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			got, err := CreateStartReceptionHandlerRequest(req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, *got)
		})
	}
}

func TestCreateAddProductHandlerRequest(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantErr  bool
		expected AddProductHandlerRequest
	}{
		{
			name:    "Valid request",
			body:    `{"type": "электроника", "pvzId": "123"}`,
			wantErr: false,
			expected: AddProductHandlerRequest{
				Type:  "электроника",
				PvzId: "123",
			},
		},
		{
			name:    "Missing type",
			body:    `{"pvzId": "123"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/products", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			got, err := CreateAddProductHandlerRequest(req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, *got)
		})
	}
}

func TestAuthError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "Not found error",
			err:      repository.NotFound,
			expected: http.StatusUnauthorized,
		},
		{
			name:     "Other error",
			err:      errors.New("other error"),
			expected: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			AuthError(rec, tt.err, nil)
			assert.Equal(t, tt.expected, rec.Code)
		})
	}
}

func TestCreateListPvzHandlerRequest(t *testing.T) {
	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		query    string
		wantErr  bool
		expected ListPvzHandlerRequest
	}{
		{
			name:    "Valid params",
			query:   "?startDate=2023-01-01T00:00:00Z&endDate=2023-01-31T00:00:00Z&page=2&limit=20",
			wantErr: false,
			expected: ListPvzHandlerRequest{
				StartDate: &startDate,
				EndDate:   &endDate,
				Page:      2,
				Limit:     20,
			},
		},
		{
			name:    "Default values",
			query:   "",
			wantErr: false,
			expected: ListPvzHandlerRequest{
				Page:  1,
				Limit: 10,
			},
		},
		{
			name:    "Invalid date format",
			query:   "?startDate=invalid_date",
			wantErr: true,
		},
		{
			name:    "Only start date invalid",
			query:   "?startDate=invalid&endDate=2023-01-31T00:00:00Z",
			wantErr: true,
		},
		{
			name:    "Only end date invalid",
			query:   "?startDate=2023-01-01T00:00:00Z&endDate=invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/pvz"+tt.query, nil)
			got, err := CreateListPvzHandlerRequest(req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for query: %s", tt.query)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.expected.StartDate != nil {
				if got.StartDate == nil {
					t.Error("Expected StartDate, got nil")
				} else if !tt.expected.StartDate.Equal(*got.StartDate) {
					t.Errorf("StartDate mismatch: expected %v, got %v", *tt.expected.StartDate, *got.StartDate)
				}
			}

			if tt.expected.EndDate != nil {
				if got.EndDate == nil {
					t.Error("Expected EndDate, got nil")
				} else if !tt.expected.EndDate.Equal(*got.EndDate) {
					t.Errorf("EndDate mismatch: expected %v, got %v", *tt.expected.EndDate, *got.EndDate)
				}
			}

			if got.Page != tt.expected.Page {
				t.Errorf("Page mismatch: expected %d, got %d", tt.expected.Page, got.Page)
			}

			if got.Limit != tt.expected.Limit {
				t.Errorf("Limit mismatch: expected %d, got %d", tt.expected.Limit, got.Limit)
			}
		})
	}
}
