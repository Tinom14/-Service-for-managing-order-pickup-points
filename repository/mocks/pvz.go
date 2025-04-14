package mocks

import (
	"avito_test/domain"
	"avito_test/usecases"
	"github.com/stretchr/testify/mock"
	"time"
)

type Pvz struct {
	mock.Mock
}

func (m *Pvz) OpenPvz(city string) (domain.Pvz, error) {
	args := m.Called(city)
	return args.Get(0).(domain.Pvz), args.Error(1)
}

func (m *Pvz) GetPvz(pvzId int) (domain.Pvz, error) {
	args := m.Called(pvzId)
	return args.Get(0).(domain.Pvz), args.Error(1)
}

func (m *Pvz) GetPvzListWithFilter(startDate, endDate *time.Time, offset, limit int) ([]usecases.PvzWithReceptions, error) {
	return []usecases.PvzWithReceptions{}, nil

}
