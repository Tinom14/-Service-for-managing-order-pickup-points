package http_test

import (
	http2 "avito_test/api/http"
	"avito_test/domain"
	"avito_test/pkg/testutils"
	"avito_test/usecases"
	"avito_test/usecases/mocks"
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReceptionHandler_StartReception(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		mockSetup    func(*mocks.Reception)
		expectedCode int
	}{
		{
			name:        "Success start reception",
			requestBody: `{"pvzId": "1"}`,
			mockSetup: func(m *mocks.Reception) {
				m.On("StartReception", 1).Return(testutils.MockReception(), nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:        "Unclosed reception exists",
			requestBody: `{"pvzId": "1"}`,
			mockSetup: func(m *mocks.Reception) {
				m.On("StartReception", 1).Return(domain.Reception{}, usecases.ErrUnclosedReception)
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.Reception)
			tt.mockSetup(mockService)
			handler := http2.NewReceptionHandler(mockService)

			req := httptest.NewRequest("POST", "/receptions", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/receptions", handler.StartReceptionHandler)
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestReceptionHandler_CloseReception(t *testing.T) {
	tests := []struct {
		name         string
		pvzId        string
		mockSetup    func(*mocks.Reception)
		expectedCode int
	}{
		{
			name:  "Success close reception",
			pvzId: "1",
			mockSetup: func(m *mocks.Reception) {
				m.On("CloseReception", 1).Return(testutils.MockReception(), nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:  "Already closed",
			pvzId: "1",
			mockSetup: func(m *mocks.Reception) {
				m.On("CloseReception", 1).Return(domain.Reception{}, usecases.ErrAlreadyClosed)
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.Reception)
			tt.mockSetup(mockService)
			handler := http2.NewReceptionHandler(mockService)

			req := httptest.NewRequest("POST", "/pvz/"+tt.pvzId+"/close_last_reception", nil)
			rec := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/pvz/{pvzId}/close_last_reception", handler.CloseReceptionHandler)
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}
