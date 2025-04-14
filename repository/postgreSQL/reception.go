package postgreSQL

import (
	"avito_test/domain"
	"avito_test/pkg/postgres_connect"
	"time"
)

type ReceptionRepo struct {
	receptions *postgres_connect.PostgresStorage
}

func NewReceptionRepo(receptions *postgres_connect.PostgresStorage) *ReceptionRepo {
	return &ReceptionRepo{receptions: receptions}
}

func (r *ReceptionRepo) StartReception(pvzId int) (domain.Reception, error) {
	now := time.Now()
	status := "in_progress"
	var id int
	err := r.receptions.Db.QueryRow(
		`INSERT INTO receptions (pvz_id, created_at, status) VALUES ($1, $2, $3) RETURNING id`,
		pvzId, now, status,
	).Scan(&id)
	if err != nil {
		return domain.Reception{}, err
	}

	return domain.Reception{Id: id, PvzId: pvzId, StartDate: now, Status: status}, nil
}

func (r *ReceptionRepo) CloseReception(pvzId int) (domain.Reception, error) {
	var reception domain.Reception

	err := r.receptions.Db.QueryRow(`
		SELECT id FROM receptions
		WHERE pvz_id = $1
		ORDER BY created_at DESC LIMIT 1`, pvzId).Scan(&reception.Id)

	if err != nil {
		return domain.Reception{}, err
	}

	_, err = r.receptions.Db.Exec(`
		UPDATE receptions SET status = 'closed'
		WHERE id = $1`, reception.Id)

	if err != nil {
		return domain.Reception{}, err
	}

	reception.Status = "closed"
	reception.PvzId = pvzId

	return reception, nil
}

func (r *ReceptionRepo) GetLastReception(pvzId int) (domain.Reception, error) {
	var rec domain.Reception
	err := r.receptions.Db.QueryRow(`
		SELECT id, created_at,status
		FROM receptions
		WHERE pvz_id = $1
		ORDER BY created_at DESC LIMIT 1`, pvzId).
		Scan(&rec.Id, &rec.StartDate, &rec.Status)

	if err != nil {
		return domain.Reception{}, err
	}

	rec.PvzId = pvzId
	return rec, nil
}

func (r *ReceptionRepo) AddProduct(pvzId int, productId int) error {
	rec, err := r.GetLastReception(pvzId)
	if err != nil {
		return err
	}

	_, err = r.receptions.Db.Exec(
		`INSERT INTO reception_products (reception_id, product_id) VALUES ($1, $2)`,
		rec.Id, productId,
	)
	return err
}

func (r *ReceptionRepo) DeleteProduct(pvzId int) (string, error) {
	rec, err := r.GetLastReception(pvzId)
	if err != nil {
		return "", err
	}

	var productId string
	err = r.receptions.Db.QueryRow(`
		SELECT product_id FROM reception_products
		WHERE reception_id = $1
		ORDER BY product_id DESC LIMIT 1`, rec.Id).Scan(&productId)

	if err != nil {
		return "", err
	}

	_, err = r.receptions.Db.Exec(`
		DELETE FROM reception_products WHERE reception_id = $1 AND product_id = $2`,
		rec.Id, productId,
	)
	if err != nil {
		return "", err
	}

	return productId, nil
}
