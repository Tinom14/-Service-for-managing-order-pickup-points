package postgreSQL

import (
	"avito_test/domain"
	"avito_test/pkg/postgres_connect"
	"time"
)

type ProductRepo struct {
	products *postgres_connect.PostgresStorage
}

func NewProductRepo(products *postgres_connect.PostgresStorage) *ProductRepo {
	return &ProductRepo{products: products}
}

func (r *ProductRepo) AddProduct(productType string) (domain.Product, error) {
	now := time.Now()
	var id int
	err := r.products.Db.QueryRow(
		`INSERT INTO products (type, added_at) VALUES ($1, $2) RETURNING id`,
		productType, now,
	).Scan(&id)
	if err != nil {
		return domain.Product{}, err
	}

	return domain.Product{
		Id:       id,
		Type:     productType,
		DateTime: now,
	}, nil
}

func (r *ProductRepo) DeleteProduct(productId int) error {
	_, err := r.products.Db.Exec(`DELETE FROM products WHERE id = $1`, productId)
	return err
}
