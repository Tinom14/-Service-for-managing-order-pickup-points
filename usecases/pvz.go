package usecases

import (
	"avito_test/domain"
	"time"
)

type PvzWithReceptions struct {
	Pvz        domain.Pvz
	Receptions []domain.ReceptionWithProducts
}

type Pvz interface {
	OpenPvz(city string) (domain.Pvz, error)
	GetPvz(pvzId int) (domain.Pvz, error)
	GetPvzListWithFilter(startDate, endDate *time.Time, page, limit int) ([]PvzWithReceptions, error)
}
