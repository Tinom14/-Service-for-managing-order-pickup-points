package service

import (
	"avito_test/domain"
	"avito_test/repository/mocks"
	"avito_test/usecases/service"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPvzService_OpenPvz(t *testing.T) {
	tests := []struct {
		name    string
		city    string
		mockPvz domain.Pvz
		mockErr error
		wantErr bool
	}{
		{
			name:    "successful open pvz",
			city:    "Moscow",
			mockPvz: domain.Pvz{Id: 1, City: "Moscow"},
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "open pvz error",
			city:    "Moscow",
			mockPvz: domain.Pvz{},
			mockErr: errors.New("repository error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.Pvz)
			mockRepo.On("OpenPvz", tt.city).Return(tt.mockPvz, tt.mockErr)

			pvzService := service.NewPvzService(mockRepo)
			pvz, err := pvzService.OpenPvz(tt.city)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockPvz, pvz)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPvzService_GetPvz(t *testing.T) {
	tests := []struct {
		name    string
		pvzId   int
		mockPvz domain.Pvz
		mockErr error
		wantErr bool
	}{
		{
			name:    "successful get pvz",
			pvzId:   1,
			mockPvz: domain.Pvz{Id: 1, City: "Moscow"},
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "get pvz error",
			pvzId:   1,
			mockPvz: domain.Pvz{},
			mockErr: errors.New("not found"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.Pvz)
			mockRepo.On("GetPvz", tt.pvzId).Return(tt.mockPvz, tt.mockErr)

			pvzService := service.NewPvzService(mockRepo)
			pvz, err := pvzService.GetPvz(tt.pvzId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockPvz, pvz)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
