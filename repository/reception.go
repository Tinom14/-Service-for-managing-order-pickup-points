package repository

import "avito_test/domain"

type Reception interface {
	StartReception(pvzId int) (domain.Reception, error)
	CloseReception(receptionId int) (domain.Reception, error)
	GetLastReception(pvzId int) (domain.Reception, error)
	AddProduct(pvzId int, productId int) error
	DeleteProduct(pvzId int) (string, error)
}
