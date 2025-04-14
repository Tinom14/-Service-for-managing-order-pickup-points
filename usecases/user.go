package usecases

import "avito_test/domain"

type User interface {
	GetToken(id string, role string) (string, error)
	Register(email string, password string, role string) (domain.User, error)
	Login(email string, password string) (string, error)
}
