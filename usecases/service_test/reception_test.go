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

func TestReceptionService_StartReception(t *testing.T) {
	tests := []struct {
		name                  string
		pvzId                 int
		mockPvz               domain.Pvz
		mockPvzErr            error
		mockReception         domain.Reception
		mockReceptionErr      error
		mockStartReception    domain.Reception
		mockStartReceptionErr error
		wantErr               bool
		expectedErr           error
	}{
		{
			name:                  "successful start reception",
			pvzId:                 1,
			mockPvz:               domain.Pvz{Id: 1},
			mockPvzErr:            nil,
			mockReception:         domain.Reception{Status: "closed"},
			mockReceptionErr:      nil,
			mockStartReception:    domain.Reception{Id: 1, PvzId: 1, Status: "in_progress"},
			mockStartReceptionErr: nil,
			wantErr:               false,
		},
		{
			name:        "pvz not found",
			pvzId:       1,
			mockPvz:     domain.Pvz{},
			mockPvzErr:  repository.NotFound,
			wantErr:     true,
			expectedErr: repository.NotFound,
		},
		{
			name:             "unclosed reception",
			pvzId:            1,
			mockPvz:          domain.Pvz{Id: 1},
			mockPvzErr:       nil,
			mockReception:    domain.Reception{Status: "in_progress"},
			mockReceptionErr: nil,
			wantErr:          true,
			expectedErr:      usecases.ErrUnclosedReception,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReceptionRepo := new(mocks.Reception)
			mockPvzRepo := new(mocks.Pvz)

			mockPvzRepo.On("GetPvz", tt.pvzId).Return(tt.mockPvz, tt.mockPvzErr)
			if tt.mockPvzErr == nil {
				mockReceptionRepo.On("GetLastReception", tt.pvzId).Return(tt.mockReception, tt.mockReceptionErr)
				if tt.mockReception.Status == "closed" {
					mockReceptionRepo.On("StartReception", tt.pvzId).Return(tt.mockStartReception, tt.mockStartReceptionErr)
				}
			}

			receptionService := service.NewReceptionService(mockReceptionRepo, mockPvzRepo)
			_, err := receptionService.StartReception(tt.pvzId)

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
		})
	}
}

func TestReceptionService_CloseReception(t *testing.T) {
	tests := []struct {
		name                  string
		pvzId                 int
		mockPvz               domain.Pvz
		mockPvzErr            error
		mockReception         domain.Reception
		mockReceptionErr      error
		mockCloseReception    domain.Reception
		mockCloseReceptionErr error
		wantErr               bool
		expectedErr           error
	}{
		{
			name:                  "successful close reception",
			pvzId:                 1,
			mockPvz:               domain.Pvz{Id: 1},
			mockPvzErr:            nil,
			mockReception:         domain.Reception{Status: "in_progress"},
			mockReceptionErr:      nil,
			mockCloseReception:    domain.Reception{Id: 1, PvzId: 1, Status: "closed"},
			mockCloseReceptionErr: nil,
			wantErr:               false,
		},
		{
			name:             "already closed",
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
			mockReceptionRepo := new(mocks.Reception)
			mockPvzRepo := new(mocks.Pvz)

			mockPvzRepo.On("GetPvz", tt.pvzId).Return(tt.mockPvz, tt.mockPvzErr)
			if tt.mockPvzErr == nil {
				mockReceptionRepo.On("GetLastReception", tt.pvzId).Return(tt.mockReception, tt.mockReceptionErr)
				if tt.mockReception.Status == "in_progress" {
					mockReceptionRepo.On("CloseReception", tt.pvzId).Return(tt.mockCloseReception, tt.mockCloseReceptionErr)
				}
			}

			receptionService := service.NewReceptionService(mockReceptionRepo, mockPvzRepo)
			reception, err := receptionService.CloseReception(tt.pvzId)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockCloseReception, reception)
			}

			mockPvzRepo.AssertExpectations(t)
			mockReceptionRepo.AssertExpectations(t)
		})
	}
}
