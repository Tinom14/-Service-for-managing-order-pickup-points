package postgreSQL

import (
	"avito_test/domain"
	"avito_test/pkg/postgres_connect"
	"avito_test/repository"
	"avito_test/usecases"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"strconv"
	"strings"
	"time"
)

type PvzRepo struct {
	pvz *postgres_connect.PostgresStorage
}

func NewPvzRepo(pvz *postgres_connect.PostgresStorage) *PvzRepo {
	return &PvzRepo{pvz: pvz}
}

func (p *PvzRepo) OpenPvz(city string) (domain.Pvz, error) {
	now := time.Now()
	var id int
	err := p.pvz.Db.QueryRow(
		`INSERT INTO pvz (city, registration_date) VALUES ($1, $2) RETURNING id`,
		city, now,
	).Scan(&id)
	if err != nil {
		return domain.Pvz{}, err
	}

	return domain.Pvz{Id: id, City: city, RegistrationDate: now}, nil
}

func (p *PvzRepo) GetPvz(pvzID int) (domain.Pvz, error) {
	row := p.pvz.Db.QueryRow(`SELECT id, city, registration_date FROM pvz WHERE id = $1`, pvzID)

	var pvz domain.Pvz
	err := row.Scan(&pvz.Id, &pvz.City, &pvz.RegistrationDate)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Pvz{}, repository.NotFound
	} else if err != nil {
		return domain.Pvz{}, err
	}
	return pvz, nil
}

func (p *PvzRepo) GetPvzListWithFilter(startDate, endDate *time.Time, offset, limit int) ([]usecases.PvzWithReceptions, error) {
	query := `
        SELECT p.id, p.city, p.registration_date, 
               r.id, r.created_at, r.status
        FROM pvz p
        LEFT JOIN receptions r ON p.id = r.pvz_id
    `

	var args []interface{}
	var where []string

	if startDate != nil {
		where = append(where, "r.created_at >= $"+strconv.Itoa(len(args)+1))
		args = append(args, *startDate)
	}
	if endDate != nil {
		where = append(where, "r.created_at <= $"+strconv.Itoa(len(args)+1))
		args = append(args, *endDate)
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	query += " ORDER BY p.id, r.created_at DESC LIMIT $" + strconv.Itoa(len(args)+1)
	args = append(args, limit)
	query += " OFFSET $" + strconv.Itoa(len(args)+1)
	args = append(args, offset)

	rows, err := p.pvz.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]*usecases.PvzWithReceptions)
	receptionIDs := make([]int, 0)

	for rows.Next() {
		var pvz domain.Pvz
		var reception domain.Reception

		err := rows.Scan(
			&pvz.Id, &pvz.City, &pvz.RegistrationDate,
			&reception.Id, &reception.StartDate, &reception.Status,
		)
		if err != nil {
			return nil, err
		}

		if _, ok := result[pvz.Id]; !ok {
			result[pvz.Id] = &usecases.PvzWithReceptions{
				Pvz:        pvz,
				Receptions: []domain.ReceptionWithProducts{},
			}
		}

		if reception.Id != 0 {
			result[pvz.Id].Receptions = append(result[pvz.Id].Receptions, domain.ReceptionWithProducts{
				Reception: reception,
				Products:  nil, // временно nil
			})
			receptionIDs = append(receptionIDs, reception.Id)
		}
	}

	productsMap, err := p.getProductsForReceptionsMap(receptionIDs)
	if err != nil {
		return nil, err
	}

	for _, pvzData := range result {
		for i := range pvzData.Receptions {
			rid := pvzData.Receptions[i].Reception.Id
			pvzData.Receptions[i].Products = productsMap[rid]
		}
	}

	var finalResult []usecases.PvzWithReceptions
	for _, v := range result {
		finalResult = append(finalResult, *v)
	}

	return finalResult, nil
}

func (p *PvzRepo) getProductsForReceptionsMap(receptionIDs []int) (map[int][]domain.Product, error) {
	if len(receptionIDs) == 0 {
		return make(map[int][]domain.Product), nil
	}

	args := make([]interface{}, len(receptionIDs))
	for i, id := range receptionIDs {
		args[i] = id
	}

	query := `
        SELECT p.id, p.type, p.added_at, rp.reception_id
        FROM products p
        JOIN reception_products rp ON p.id = rp.product_id
        WHERE rp.reception_id = ANY($1)
    `

	rows, err := p.pvz.Db.Query(query, pq.Array(receptionIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int][]domain.Product)

	for rows.Next() {
		var product domain.Product
		var receptionID int
		if err := rows.Scan(&product.Id, &product.Type, &product.DateTime, &receptionID); err != nil {
			return nil, err
		}
		result[receptionID] = append(result[receptionID], product)
	}

	return result, nil
}
