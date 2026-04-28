package dto

import (
	"payment-processor/internal/payment/adapters/outbound/service"
	"payment-processor/internal/payment/domain"
	"time"
)

const (
	StatusPending  = string(domain.StatusPending)
	StatusApproved = string(domain.StatusApproved)
	StatusRejected = string(domain.StatusRejected)
)

type CreatePaymentInput struct {
	Item    Item   `json:"item"`
	Payer   Payer  `json:"payer"`
	BackURL string `json:"back_url"`
}

type Item struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	CurrencyID string  `json:"currency_id"`
}

type Payer struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
}

type CreatePaymentOutput struct {
	ID                string    `json:"id"`
	InitPoint         string    `json:"init_point"`
	SandboxInitPoint  string    `json:"sandbox_init_point"`
	ExternalReference string    `json:"external_reference"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func ToPayment(prot, name, surname, email string, qtdy int, total, stotal float64, desc, prefId, extRef, iPoint, sIPoint string) (*domain.Payment, error) {
	return domain.NewPayment(
		prot,
		name,
		surname,
		email,
		qtdy,
		total,
		stotal,
		desc,
		prefId,
		extRef,
		iPoint,
		sIPoint,
	)
}

func FromPayment(payment *domain.Payment) *CreatePaymentOutput {
	return &CreatePaymentOutput{
		ID:                payment.ID,
		InitPoint:         payment.InitPoint,
		SandboxInitPoint:  payment.SandboxInitPoint,
		ExternalReference: payment.ExternalReference,
		CreatedAt:         payment.CreatedAt,
		UpdatedAt:         payment.UpdatedAt,
	}
}

func FromPayments(payments []*domain.Payment) *[]CreatePaymentOutput {
	outputs := make([]CreatePaymentOutput, len(payments))
	for i, payment := range payments {
		outputs[i] = *FromPayment(payment)
	}
	return &outputs
}

func MapItem(dtoItem Item) service.Item {
	return service.Item{
		ID:         dtoItem.ID,
		Title:      dtoItem.Title,
		Quantity:   dtoItem.Quantity,
		UnitPrice:  dtoItem.UnitPrice,
		CurrencyID: dtoItem.CurrencyID,
	}
}

func MapPayer(dtoPayer Payer) service.Payer {
	return service.Payer{
		Name:        dtoPayer.Name,
		Surname:     dtoPayer.Surname,
		Email:       dtoPayer.Email,
		DateCreated: dtoPayer.DateCreated,
	}
}
