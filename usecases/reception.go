package usecases

import "avito_test/domain"

type Reception interface {
	StartReception(pvzId int) (domain.Reception, error)
	CloseReception(pvzId int) (domain.Reception, error)
	CheckPvz(pvzId int) error
}
