package service

import (
	"avito_test/domain"
	"avito_test/repository"
	"avito_test/usecases"
	"strconv"
)

type Product struct {
	productRepo   repository.Product
	receptionRepo repository.Reception
	pvzRepo       repository.Pvz
}

func NewProductService(productRepo repository.Product, receptionRepo repository.Reception, pvzRepo repository.Pvz) *Product {
	return &Product{productRepo: productRepo, receptionRepo: receptionRepo, pvzRepo: pvzRepo}
}

func (p *Product) AddProduct(sort string, pvzId int) (domain.Product, error) {
	if _, err := p.pvzRepo.GetPvz(pvzId); err != nil {
		return domain.Product{}, err
	}

	lastReception, err := p.receptionRepo.GetLastReception(pvzId)
	if err != nil {
		return domain.Product{}, repository.NotFound
	}
	lastReceptionStatus := lastReception.Status
	if lastReceptionStatus == "closed" {
		return domain.Product{}, usecases.ErrAlreadyClosed
	}

	product, err := p.productRepo.AddProduct(sort)
	if err != nil {
		return domain.Product{}, err
	}
	_ = p.receptionRepo.AddProduct(pvzId, product.Id)
	return product, err
}

func (p *Product) DeleteProduct(pvzId int) error {
	if _, err := p.pvzRepo.GetPvz(pvzId); err != nil {
		return repository.NotFound
	}

	lastReception, err := p.receptionRepo.GetLastReception(pvzId)
	if err != nil {
		return repository.NotFound
	}
	if lastReception.Status == "closed" {
		return usecases.ErrAlreadyClosed
	}
	productId, err := p.receptionRepo.DeleteProduct(pvzId)
	if err != nil {
		return err
	}
	productIdInt, _ := strconv.Atoi(productId)
	err = p.productRepo.DeleteProduct(productIdInt)
	return err
}
