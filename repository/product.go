package repository

import "avito_test/domain"

type Product interface {
	AddProduct(sort string) (product domain.Product, err error)
	DeleteProduct(productId int) (err error)
}
