package mocks

import (
	"avito_test/domain"
	"github.com/stretchr/testify/mock"
)

type User struct {
	mock.Mock
}

func (m *User) GetToken(id string, role string) (string, error) {
	args := m.Called(id, role)
	return args.String(0), args.Error(1)
}

func (m *User) Register(email string, password string, role string) (domain.User, error) {
	args := m.Called(email, password, role)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *User) Login(email string, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}
