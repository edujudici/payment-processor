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
	ID                string
	Protocol          string
	Username          string
	Surname           string
	Email             string
	Status            Status
	Quantity          int
	Total             float64
	Subtotal          float64
	Description       string
	PreferenceID      string
	ExternalReference string
	InitPoint         string
	SandboxInitPoint  string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewPayment(prot, name, surname, email string, qtdy int, total, stotal float64, desc, prefId, extRef, iPoint, sIPoint string) (*Payment, error) {

	if qtdy <= 0 {
		return nil, ErrInvalidQuantity
	}

	return &Payment{
		ID:                uuid.New().String(),
		Protocol:          prot,
		Username:          name,
		Surname:           surname,
		Email:             email,
		Status:            StatusPending,
		Quantity:          qtdy,
		Total:             total,
		Subtotal:          stotal,
		Description:       desc,
		PreferenceID:      prefId,
		ExternalReference: extRef,
		InitPoint:         iPoint,
		SandboxInitPoint:  sIPoint,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
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
