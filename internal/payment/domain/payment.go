package domain

import "time"

type Payment struct {
	ID        string
	Amount    float64
	Status    string
	CreatedAt time.Time
}
