package repository_test

import (
	"avito_test/domain"
	"avito_test/pkg/postgres_connect"
	"avito_test/repository/postgreSQL"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestProductRepo_AddProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgreSQL.NewProductRepo(&postgres_connect.PostgresStorage{Db: db})

	tests := []struct {
		name        string
		productType string
		mock        func()
		want        domain.Product
		wantErr     bool
	}{
		{
			name:        "success",
			productType: "apple",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`INSERT INTO products`).
					WithArgs("apple", sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			want: domain.Product{
				Id:       1,
				Type:     "apple",
				DateTime: time.Now(),
			},
			wantErr: false,
		},
		{
			name:        "database error",
			productType: "apple",
			mock: func() {
				mock.ExpectQuery(`INSERT INTO products`).
					WithArgs("apple", sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			want:    domain.Product{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repo.AddProduct(tt.productType)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.Type, got.Type)
				assert.NotZero(t, got.DateTime)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProductRepo_DeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgreSQL.NewProductRepo(&postgres_connect.PostgresStorage{Db: db})

	tests := []struct {
		name      string
		productId int
		mock      func()
		wantErr   bool
	}{
		{
			name:      "success",
			productId: 1,
			mock: func() {
				mock.ExpectExec(`DELETE FROM products`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:      "database error",
			productId: 1,
			mock: func() {
				mock.ExpectExec(`DELETE FROM products`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.DeleteProduct(tt.productId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProductRepo_DeleteProduct_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := postgreSQL.NewProductRepo(&postgres_connect.PostgresStorage{Db: db})

	mock.ExpectExec(`DELETE FROM products`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.DeleteProduct(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepo_AddProduct_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := postgreSQL.NewProductRepo(&postgres_connect.PostgresStorage{Db: db})

	mock.ExpectQuery(`INSERT INTO products`).
		WithArgs("apple", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("invalid")) // Неправильный тип для id

	_, err = repo.AddProduct("apple")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
