package repository

import "avito_test/domain"

type User interface {
	Register(email string, password string, role string) (domain.User, error)
	Login(email string) (domain.User, error)
}
