package repository

import (
	"avito_test/domain"
	"avito_test/pkg/postgres_connect"
	"avito_test/repository"
	"avito_test/repository/postgreSQL"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepo_Register(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgreSQL.NewUserRepo(&postgres_connect.PostgresStorage{Db: db})

	tests := []struct {
		name     string
		email    string
		password string
		role     string
		mock     func()
		want     domain.User
		wantErr  bool
	}{
		{
			name:     "success",
			email:    "test@example.com",
			password: "password123",
			role:     "user",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs("test@example.com", "password123", "user").
					WillReturnRows(rows)
			},
			want: domain.User{
				Id:       1,
				Email:    "test@example.com",
				Password: "password123",
				Role:     "user",
			},
			wantErr: false,
		},
		{
			name:     "email already exists",
			email:    "test@example.com",
			password: "password123",
			role:     "user",
			mock: func() {
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs("test@example.com", "password123", "user").
					WillReturnError(&pq.Error{Code: "23505"})
			},
			want:    domain.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.Register(tt.email, tt.password, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.name == "email already exists" {
					assert.ErrorIs(t, err, repository.ErrEmailAlreadyExists)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.Password, got.Password)
				assert.Equal(t, tt.want.Role, got.Role)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepo_Login(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgreSQL.NewUserRepo(&postgres_connect.PostgresStorage{Db: db})

	tests := []struct {
		name    string
		email   string
		mock    func()
		want    domain.User
		wantErr bool
	}{
		{
			name:  "success",
			email: "test@example.com",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "role"}).
					AddRow(1, "test@example.com", "hashedpassword", "user")
				mock.ExpectQuery(`SELECT id, email, password_hash, role`).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			want: domain.User{
				Id:       1,
				Email:    "test@example.com",
				Password: "hashedpassword",
				Role:     "user",
			},
			wantErr: false,
		},
		{
			name:  "not found",
			email: "test@example.com",
			mock: func() {
				mock.ExpectQuery(`SELECT id, email, password_hash, role`).
					WithArgs("test@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			want:    domain.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.Login(tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				if errors.Is(err, sql.ErrNoRows) {
					assert.ErrorIs(t, err, repository.NotFound)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.Password, got.Password)
				assert.Equal(t, tt.want.Role, got.Role)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepo_Register_UnknownDBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := postgreSQL.NewUserRepo(&postgres_connect.PostgresStorage{Db: db})

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs("test@example.com", "password123", "user").
		WillReturnError(errors.New("unknown error"))

	_, err = repo.Register("test@example.com", "password123", "user")
	assert.Error(t, err)
	assert.NotErrorIs(t, err, repository.ErrEmailAlreadyExists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_Login_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := postgreSQL.NewUserRepo(&postgres_connect.PostgresStorage{Db: db})

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "role"}).
		AddRow("invalid", "test@example.com", "hashedpassword", "user") // Неправильный тип для id
	mock.ExpectQuery(`SELECT id, email, password_hash, role`).
		WithArgs("test@example.com").
		WillReturnRows(rows)

	_, err = repo.Login("test@example.com")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
