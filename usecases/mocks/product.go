package mocks

import (
	"avito_test/domain"
	"github.com/stretchr/testify/mock"
)

type Product struct {
	mock.Mock
}

func (m *Product) AddProduct(sort string, pvzId int) (domain.Product, error) {
	args := m.Called(sort, pvzId)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *Product) DeleteProduct(pvzId int) error {
	args := m.Called(pvzId)
	return args.Error(0)
}
