package mocks

import (
	"avito_test/domain"
	"github.com/stretchr/testify/mock"
)

type Reception struct {
	mock.Mock
}

func (m *Reception) StartReception(pvzId int) (domain.Reception, error) {
	args := m.Called(pvzId)
	return args.Get(0).(domain.Reception), args.Error(1)
}

func (m *Reception) CloseReception(pvzId int) (domain.Reception, error) {
	args := m.Called(pvzId)
	return args.Get(0).(domain.Reception), args.Error(1)
}

func (m *Reception) GetLastReception(pvzId int) (domain.Reception, error) {
	args := m.Called(pvzId)
	return args.Get(0).(domain.Reception), args.Error(1)
}

func (m *Reception) AddProduct(receptionId int, productId int) error {
	args := m.Called(receptionId, productId)
	return args.Error(0)
}

func (m *Reception) DeleteProduct(pvzId int) (string, error) {
	args := m.Called(pvzId)
	return args.String(0), args.Error(1)
}
