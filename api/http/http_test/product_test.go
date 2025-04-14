package http_test

import (
	http2 "avito_test/api/http"
	"avito_test/domain"
	"avito_test/pkg/testutils"
	"avito_test/repository"
	"avito_test/usecases"
	"avito_test/usecases/mocks"
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProductHandler_AddProduct(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		mockSetup    func(*mocks.Product)
		expectedCode int
	}{
		{
			name:        "Success add product",
			requestBody: `{"type": "электроника", "pvzId": "1"}`,
			mockSetup: func(m *mocks.Product) {
				m.On("AddProduct", "электроника", 1).Return(testutils.MockProduct(), nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:        "Reception closed",
			requestBody: `{"type": "электроника", "pvzId": "1"}`,
			mockSetup: func(m *mocks.Product) {
				m.On("AddProduct", "электроника", 1).Return(domain.Product{}, usecases.ErrAlreadyClosed)
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.Product)
			tt.mockSetup(mockService)
			handler := http2.NewProductHandler(mockService)

			req := httptest.NewRequest("POST", "/products", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/products", handler.AddProductHandler)
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	tests := []struct {
		name         string
		pvzId        string
		mockSetup    func(*mocks.Product)
		expectedCode int
	}{
		{
			name:  "Success delete product",
			pvzId: "1",
			mockSetup: func(m *mocks.Product) {
				m.On("DeleteProduct", 1).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:  "No products to delete",
			pvzId: "1",
			mockSetup: func(m *mocks.Product) {
				m.On("DeleteProduct", 1).Return(repository.NotFound)
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.Product)
			tt.mockSetup(mockService)
			handler := http2.NewProductHandler(mockService)

			req := httptest.NewRequest("POST", "/pvz/"+tt.pvzId+"/delete_last_product", nil)
			rec := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/pvz/{pvzId}/delete_last_product", handler.DeleteProductHandler)
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestProductHandler_AddProduct_InvalidJSON(t *testing.T) {
	handler := &http2.Product{Service: new(mocks.Product)}

	req := httptest.NewRequest("POST", "/products", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.AddProductHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProductHandler_DeleteProduct_InvalidPVZId(t *testing.T) {
	handler := &http2.Product{Service: new(mocks.Product)}

	req := httptest.NewRequest("POST", "/pvz/invalid/delete_last_product", nil)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/pvz/{pvzId}/delete_last_product", handler.DeleteProductHandler)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
