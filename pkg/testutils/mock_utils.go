package testutils

import (
	"avito_test/domain"
	"time"
)

func MockPvz() domain.Pvz {
	return domain.Pvz{
		Id:               1,
		City:             "Москва",
		RegistrationDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func MockReception() domain.Reception {
	return domain.Reception{
		Id:        1,
		PvzId:     1,
		Status:    "in_progress",
		StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func MockProduct() domain.Product {
	return domain.Product{
		Id:       1,
		Type:     "электроника",
		DateTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}
