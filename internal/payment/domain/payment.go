package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

type Payment struct {
	ID           string
	PreferenceID string
	Amount       float64
	Status       Status
	Description  string
	PaymentType  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewPayment(preferenceId string, amount float64, description, paymentType string) (*Payment, error) {

	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	return &Payment{
		ID:           uuid.New().String(),
		PreferenceID: preferenceId,
		Amount:       amount,
		Status:       StatusPending,
		Description:  description,
		PaymentType:  paymentType,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (p *Payment) UpdateStatus(newStatus Status) error {
	if p.Status != StatusPending {
		return ErrInvalidStatus
	}

	p.Status = newStatus
	p.UpdatedAt = time.Now()
	return nil
}
