package postgreSQL

import (
	"avito_test/domain"
	"avito_test/pkg/postgres_connect"
	"avito_test/repository"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type UserRepo struct {
	users *postgres_connect.PostgresStorage
}

func NewUserRepo(users *postgres_connect.PostgresStorage) *UserRepo {
	return &UserRepo{users: users}
}

func (u *UserRepo) Register(email string, password string, role string) (domain.User, error) {
	var id int
	err := u.users.Db.QueryRow(
		`INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id`,
		email, password, role,
	).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return domain.User{}, repository.ErrEmailAlreadyExists
		}
		return domain.User{}, err
	}
	return domain.User{Id: id, Email: email, Password: password, Role: role}, nil
}

func (u *UserRepo) Login(email string) (domain.User, error) {
	row := u.users.Db.QueryRow(
		`SELECT id, email, password_hash, role FROM users WHERE email = $1`,
		email,
	)

	var user domain.User
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, repository.NotFound
	} else if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
