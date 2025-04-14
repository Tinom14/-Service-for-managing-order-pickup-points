package domain

import "time"

type Pvz struct {
	Id               int       `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}
