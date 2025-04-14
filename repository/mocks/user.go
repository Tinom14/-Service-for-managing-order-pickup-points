package mocks

import (
	"avito_test/domain"
	"github.com/stretchr/testify/mock"
)

type User struct {
	mock.Mock
}

func (m *User) Register(email string, password string, role string) (domain.User, error) {
	args := m.Called(email, password, role)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *User) Login(email string) (domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(domain.User), args.Error(1)
}
