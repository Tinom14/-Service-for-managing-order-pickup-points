package usecases

import "avito_test/domain"

type Product interface {
	AddProduct(sort string, pvzId int) (domain.Product, error)
	DeleteProduct(pvzId int) error
}
