package domain

import "time"

type Reception struct {
	Id        int       `json:"id"`
	StartDate time.Time `json:"startDate"`
	PvzId     int       `json:"pvzId"`
	Status    string    `json:"status"`
}
