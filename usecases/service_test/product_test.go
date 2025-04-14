package service

import (
	"avito_test/domain"
	"avito_test/repository"
	"avito_test/repository/mocks"
	"avito_test/usecases"
	"avito_test/usecases/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductService_AddProduct(t *testing.T) {
	tests := []struct {
		name             string
		sort             string
		pvzId            int
		mockPvz          domain.Pvz
		mockPvzErr       error
		mockReception    domain.Reception
		mockReceptionErr error
		mockProduct      domain.Product
		mockProductErr   error
		wantErr          bool
		expectedErr      error
	}{
		{
			name:             "successful add product",
			sort:             "apple",
			pvzId:            1,
			mockPvz:          domain.Pvz{Id: 1},
			mockPvzErr:       nil,
			mockReception:    domain.Reception{Status: "in_progress"},
			mockReceptionErr: nil,
			mockProduct:      domain.Product{Id: 1, Type: "apple"},
			mockProductErr:   nil,
			wantErr:          false,
		},
		{
			name:        "pvz not found",
			sort:        "apple",
			pvzId:       1,
			mockPvz:     domain.Pvz{},
			mockPvzErr:  repository.NotFound,
			wantErr:     true,
			expectedErr: repository.NotFound,
		},
		{
			name:             "reception not found",
			sort:             "apple",
			pvzId:            1,
			mockPvz:          domain.Pvz{Id: 1},
			mockPvzErr:       nil,
			mockReception:    domain.Reception{},
			mockReceptionErr: repository.NotFound,
			wantErr:          true,
			expectedErr:      repository.NotFound,
		},
		{
			name:             "reception already closed",
			sort:             "apple",
			pvzId:            1,
			mockPvz:          domain.Pvz{Id: 1},
			mockPvzErr:       nil,
			mockReception:    domain.Reception{Status: "closed"},
			mockReceptionErr: nil,
			wantErr:          true,
			expectedErr:      usecases.ErrAlreadyClosed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPvzRepo := new(mocks.Pvz)
			mockReceptionRepo := new(mocks.Reception)
			mockProductRepo := new(mocks.Product)

			mockPvzRepo.On("GetPvz", tt.pvzId).Return(tt.mockPvz, tt.mockPvzErr)
			if tt.mockPvzErr == nil {
				mockReceptionRepo.On("GetLastReception", tt.pvzId).Return(tt.mockReception, tt.mockReceptionErr)
			}
			if tt.mockReceptionErr == nil && tt.mockReception.Status == "in_progress" {
				mockProductRepo.On("AddProduct", tt.sort).Return(tt.mockProduct, tt.mockProductErr)
				if tt.mockProductErr == nil {
					mockReceptionRepo.On("AddProduct", tt.pvzId, tt.mockProduct.Id).Return(nil)
				}
			}

			productService := service.NewProductService(mockProductRepo, mockReceptionRepo, mockPvzRepo)
			product, err := productService.AddProduct(tt.sort, tt.pvzId)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockProduct, product)
			}

			mockPvzRepo.AssertExpectations(t)
			mockReceptionRepo.AssertExpectations(t)
			mockProductRepo.AssertExpectations(t)
		})
	}
}

func TestProductService_DeleteProduct(t *testing.T) {
	tests := []struct {
		name                 string
		pvzId                int
		mockPvz              domain.Pvz
		mockPvzErr           error
		mockReception        domain.Reception
		mockReceptionErr     error
		mockDeleteProductId  string
		mockDeleteProductErr error
		mockProductErr       error
		wantErr              bool
		expectedErr          error
	}{
		{
			name:                 "successful delete product",
			pvzId:                1,
			mockPvz:              domain.Pvz{Id: 1},
			mockPvzErr:           nil,
			mockReception:        domain.Reception{Status: "in_progress"},
			mockReceptionErr:     nil,
			mockDeleteProductId:  "1",
			mockDeleteProductErr: nil,
			mockProductErr:       nil,
			wantErr:              false,
		},
		{
			name:             "reception closed",
			pvzId:            1,
			mockPvz:          domain.Pvz{Id: 1},
			mockPvzErr:       nil,
			mockReception:    domain.Reception{Status: "closed"},
			mockReceptionErr: nil,
			wantErr:          true,
			expectedErr:      usecases.ErrAlreadyClosed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPvzRepo := new(mocks.Pvz)
			mockReceptionRepo := new(mocks.Reception)
			mockProductRepo := new(mocks.Product)

			mockPvzRepo.On("GetPvz", tt.pvzId).Return(tt.mockPvz, tt.mockPvzErr)
			if tt.mockPvzErr == nil {
				mockReceptionRepo.On("GetLastReception", tt.pvzId).Return(tt.mockReception, tt.mockReceptionErr)
			}
			if tt.mockReceptionErr == nil && tt.mockReception.Status == "in_progress" {
				mockReceptionRepo.On("DeleteProduct", tt.pvzId).Return(tt.mockDeleteProductId, tt.mockDeleteProductErr)
				if tt.mockDeleteProductErr == nil {
					mockProductRepo.On("DeleteProduct", 1).Return(tt.mockProductErr)
				}
			}

			productService := service.NewProductService(mockProductRepo, mockReceptionRepo, mockPvzRepo)
			err := productService.DeleteProduct(tt.pvzId)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}

			mockPvzRepo.AssertExpectations(t)
			mockReceptionRepo.AssertExpectations(t)
			mockProductRepo.AssertExpectations(t)
		})
	}
}
