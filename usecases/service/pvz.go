package service

import (
	"avito_test/domain"
	"avito_test/repository"
	"avito_test/usecases"
	"time"
)

type Pvz struct {
	repo repository.Pvz
}

func NewPvzService(repo repository.Pvz) *Pvz {
	return &Pvz{repo: repo}
}

func (p *Pvz) OpenPvz(city string) (domain.Pvz, error) {
	return p.repo.OpenPvz(city)
}

func (p *Pvz) GetPvz(pvzId int) (domain.Pvz, error) {
	return p.repo.GetPvz(pvzId)
}

func (p *Pvz) GetPvzListWithFilter(startDate, endDate *time.Time, page, limit int) ([]usecases.PvzWithReceptions, error) {
	offset := (page - 1) * limit
	return p.repo.GetPvzListWithFilter(startDate, endDate, offset, limit)
}
