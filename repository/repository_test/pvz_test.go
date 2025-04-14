package repository

import (
	"avito_test/domain"
	"avito_test/pkg/postgres_connect"
	"avito_test/repository"
	"avito_test/repository/postgreSQL"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPvzRepo_OpenPvz(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgreSQL.NewPvzRepo(&postgres_connect.PostgresStorage{Db: db})

	tests := []struct {
		name    string
		city    string
		mock    func()
		want    domain.Pvz
		wantErr bool
	}{
		{
			name: "success",
			city: "Moscow",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`INSERT INTO pvz`).
					WithArgs("Moscow", sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			want: domain.Pvz{
				Id:               1,
				City:             "Moscow",
				RegistrationDate: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "database error",
			city: "Moscow",
			mock: func() {
				mock.ExpectQuery(`INSERT INTO pvz`).
					WithArgs("Moscow", sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			want:    domain.Pvz{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.OpenPvz(tt.city)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.City, got.City)
				assert.NotZero(t, got.RegistrationDate)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPvzRepo_GetPvz(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgreSQL.NewPvzRepo(&postgres_connect.PostgresStorage{Db: db})

	tests := []struct {
		name    string
		pvzID   int
		mock    func()
		want    domain.Pvz
		wantErr bool
	}{
		{
			name:  "success",
			pvzID: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "city", "registration_date"}).
					AddRow(1, "Moscow", time.Now())
				mock.ExpectQuery(`SELECT id, city, registration_date`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: domain.Pvz{
				Id:               1,
				City:             "Moscow",
				RegistrationDate: time.Now(),
			},
			wantErr: false,
		},
		{
			name:  "not found",
			pvzID: 1,
			mock: func() {
				mock.ExpectQuery(`SELECT id, city, registration_date`).
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			want:    domain.Pvz{},
			wantErr: true,
		},
		{
			name:  "database error",
			pvzID: 1,
			mock: func() {
				mock.ExpectQuery(`SELECT id, city, registration_date`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			want:    domain.Pvz{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.GetPvz(tt.pvzID)

			if tt.wantErr {
				assert.Error(t, err)
				if errors.Is(err, sql.ErrNoRows) {
					assert.ErrorIs(t, err, repository.NotFound)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.City, got.City)
				assert.NotZero(t, got.RegistrationDate)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
