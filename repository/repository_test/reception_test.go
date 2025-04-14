package repository

import (
	"avito_test/domain"
	"avito_test/pkg/postgres_connect"
	"avito_test/repository/postgreSQL"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReceptionRepo_StartReception(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgreSQL.NewReceptionRepo(&postgres_connect.PostgresStorage{Db: db})

	tests := []struct {
		name    string
		pvzId   int
		mock    func()
		want    domain.Reception
		wantErr bool
	}{
		{
			name:  "success",
			pvzId: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`INSERT INTO receptions`).
					WithArgs(1, sqlmock.AnyArg(), "in_progress").
					WillReturnRows(rows)
			},
			want: domain.Reception{
				Id:        1,
				PvzId:     1,
				StartDate: time.Now(),
				Status:    "in_progress",
			},
			wantErr: false,
		},
		{
			name:  "database error",
			pvzId: 1,
			mock: func() {
				mock.ExpectQuery(`INSERT INTO receptions`).
					WithArgs(1, sqlmock.AnyArg(), "in_progress").
					WillReturnError(sql.ErrConnDone)
			},
			want:    domain.Reception{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.StartReception(tt.pvzId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.PvzId, got.PvzId)
				assert.Equal(t, tt.want.Status, got.Status)
				assert.NotZero(t, got.StartDate)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReceptionRepo_CloseReception(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgreSQL.NewReceptionRepo(&postgres_connect.PostgresStorage{Db: db})

	tests := []struct {
		name    string
		pvzId   int
		mock    func()
		want    domain.Reception
		wantErr bool
	}{
		{
			name:  "success",
			pvzId: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`SELECT id FROM receptions`).
					WithArgs(1).
					WillReturnRows(rows)
				mock.ExpectExec(`UPDATE receptions SET status`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			want: domain.Reception{
				Id:     1,
				PvzId:  1,
				Status: "closed",
			},
			wantErr: false,
		},
		{
			name:  "database error",
			pvzId: 1,
			mock: func() {
				mock.ExpectQuery(`SELECT id FROM receptions`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			want:    domain.Reception{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.CloseReception(tt.pvzId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.PvzId, got.PvzId)
				assert.Equal(t, tt.want.Status, got.Status)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReceptionRepo_GetLastReception_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := postgreSQL.NewReceptionRepo(&postgres_connect.PostgresStorage{Db: db})

	mock.ExpectQuery(`SELECT id, created_at,status`).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetLastReception(1)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReceptionRepo_AddProduct_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := postgreSQL.NewReceptionRepo(&postgres_connect.PostgresStorage{Db: db})

	rows := sqlmock.NewRows([]string{"id", "created_at", "status"}).
		AddRow(1, time.Now(), "in_progress")
	mock.ExpectQuery(`SELECT id, created_at,status`).
		WithArgs(1).
		WillReturnRows(rows)

	mock.ExpectExec(`INSERT INTO reception_products`).
		WithArgs(1, 1).
		WillReturnError(sql.ErrConnDone)

	err = repo.AddProduct(1, 1)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReceptionRepo_DeleteProduct_NoProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := postgreSQL.NewReceptionRepo(&postgres_connect.PostgresStorage{Db: db})

	rows := sqlmock.NewRows([]string{"id", "created_at", "status"}).
		AddRow(1, time.Now(), "in_progress")
	mock.ExpectQuery(`SELECT id, created_at,status`).
		WithArgs(1).
		WillReturnRows(rows)

	mock.ExpectQuery(`SELECT product_id FROM reception_products`).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.DeleteProduct(1)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
