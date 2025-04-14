package repository

import (
	"avito_test/domain"
	"avito_test/usecases"
	"time"
)

type Pvz interface {
	OpenPvz(city string) (domain.Pvz, error)
	GetPvz(pvzId int) (domain.Pvz, error)
	GetPvzListWithFilter(startDate, endDate *time.Time, offset, limit int) ([]usecases.PvzWithReceptions, error)
}
