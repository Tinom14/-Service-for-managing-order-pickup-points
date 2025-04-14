package mocks

import (
	"avito_test/domain"
	"github.com/stretchr/testify/mock"
)

type Product struct {
	mock.Mock
}

func (m *Product) AddProduct(productType string) (domain.Product, error) {
	args := m.Called(productType)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *Product) DeleteProduct(productId int) error {
	args := m.Called(productId)
	return args.Error(0)
}
