package service

import (
	"avito_test/domain"
	"avito_test/repository"
	"avito_test/usecases"
)

type Reception struct {
	repo    repository.Reception
	pvzRepo repository.Pvz
}

func NewReceptionService(receptionRepo repository.Reception, pvzRepo repository.Pvz) *Reception {
	return &Reception{
		repo:    receptionRepo,
		pvzRepo: pvzRepo,
	}
}

func (r *Reception) StartReception(pvzId int) (domain.Reception, error) {
	if err := r.CheckPvz(pvzId); err != nil {
		return domain.Reception{}, err
	}

	LastReception, _ := r.repo.GetLastReception(pvzId)
	LastReceptionStatus := LastReception.Status
	if LastReceptionStatus == "in_progress" {
		return domain.Reception{}, usecases.ErrUnclosedReception
	}

	return r.repo.StartReception(pvzId)
}

func (r *Reception) CloseReception(pvzId int) (domain.Reception, error) {
	if err := r.CheckPvz(pvzId); err != nil {
		return domain.Reception{}, err
	}

	LastReception, _ := r.repo.GetLastReception(pvzId)
	LastReceptionStatus := LastReception.Status
	if LastReceptionStatus == "closed" {
		return domain.Reception{}, usecases.ErrAlreadyClosed
	}

	return r.repo.CloseReception(pvzId)
}

func (r *Reception) CheckPvz(pvzId int) error {
	_, err := r.pvzRepo.GetPvz(pvzId)
	if err != nil {
		return err
	}
	return nil
}
