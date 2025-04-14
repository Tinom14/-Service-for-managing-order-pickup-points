package domain

import "time"

type Product struct {
	Id       int       `json:"id"`
	DateTime time.Time `json:"dateTime"`
	Type     string    `json:"type"`
}
