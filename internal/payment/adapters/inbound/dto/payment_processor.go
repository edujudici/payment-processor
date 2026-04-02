package dto

import (
	"payment-processor/internal/payment/domain"
	"time"
)

const (
	StatusPending  = string(domain.StatusPending)
	StatusApproved = string(domain.StatusApproved)
	StatusRejected = string(domain.StatusRejected)
)

type CreatePaymentInput struct {
	Amount       float64 `json:"amount"`
	PreferenceID string  `json:"preference_id"`
	Status       string  `json:"status"`
	Description  string  `json:"description"`
	PaymentType  string  `json:"payment_type"`
}

type CreatePaymentOutput struct {
	ID          string    `json:"id"`
	CheckoutURL string    `json:"checkout_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToPayment(input CreatePaymentInput) (*domain.Payment, error) {
	return domain.NewPayment(
		input.PreferenceID,
		input.Amount,
		domain.Status(input.Status),
		input.Description,
		input.PaymentType,
	)
}

func FromPayment(payment *domain.Payment) *CreatePaymentOutput {
	return &CreatePaymentOutput{
		ID:          payment.ID,
		CheckoutURL: "https://checkout.mercadopago.com.br/redirect?pref_id=" + payment.PreferenceID,
		CreatedAt:   payment.CreatedAt,
		UpdatedAt:   payment.UpdatedAt,
	}
}
