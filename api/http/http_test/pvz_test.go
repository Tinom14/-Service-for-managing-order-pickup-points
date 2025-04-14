package http_test

import (
	http2 "avito_test/api/http"
	"avito_test/domain"
	"avito_test/pkg/testutils"
	"avito_test/usecases"
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

func TestPvzHandler_OpenPvz(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		mockSetup    func(*mocks.Pvz)
		expectedCode int
	}{
		{
			name:        "Success create PVZ",
			requestBody: `{"city": "Москва"}`,
			mockSetup: func(m *mocks.Pvz) {
				m.On("OpenPvz", "Москва").Return(testutils.MockPvz(), nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Invalid city",
			requestBody:  `{"city": "Новосибирск"}`,
			mockSetup:    func(m *mocks.Pvz) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.Pvz)
			tt.mockSetup(mockService)
			handler := http2.NewPvzHandler(mockService)

			req := httptest.NewRequest("POST", "/pvz", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/pvz", handler.OpenPvzHandler)
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestPvzHandler_GetPvzList(t *testing.T) {
	mockService := new(mocks.Pvz)
	handler := http2.NewPvzHandler(mockService)

	expectedPvz := domain.Pvz{Id: 1, City: "Москва"}
	expectedReceptions := []usecases.PvzWithReceptions{
		{
			Pvz: expectedPvz,
			Receptions: []domain.ReceptionWithProducts{
				{
					Reception: domain.Reception{Id: 1},
					Products:  []domain.Product{{Id: 1}},
				},
			},
		},
	}

	mockService.On("GetPvzListWithFilter", mock.Anything, mock.Anything, 1, 10).
		Return(expectedReceptions, nil)

	req := httptest.NewRequest("GET", "/pvz?page=1&limit=10", nil)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/pvz", handler.GetPvzListHandler)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []usecases.PvzWithReceptions
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, expectedPvz.Id, response[0].Pvz.Id)

	mockService.AssertExpectations(t)
}

func TestPvzHandler_OpenPvz_DBError(t *testing.T) {
	mockService := new(mocks.Pvz)
	handler := http2.NewPvzHandler(mockService)

	mockService.On("OpenPvz", "Москва").Return(domain.Pvz{}, errors.New("db error"))

	body := `{"city": "Москва"}`
	req := httptest.NewRequest("POST", "/pvz", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	handler.OpenPvzHandler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
