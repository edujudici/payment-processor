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
	Status      string    `json:"status"`
	CheckoutURL string    `json:"checkout_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToPayment(input CreatePaymentInput) (*domain.Payment, error) {
	return domain.NewPayment(
		input.PreferenceID,
		input.Amount,
		input.Description,
		input.PaymentType,
	)
}

func FromPayment(payment *domain.Payment) *CreatePaymentOutput {
	return &CreatePaymentOutput{
		ID:          payment.ID,
		Status:      string(payment.Status),
		CheckoutURL: "https://checkout.mercadopago.com.br/redirect?pref_id=" + payment.PreferenceID,
		CreatedAt:   payment.CreatedAt,
		UpdatedAt:   payment.UpdatedAt,
	}
}

func FromPayments(payments []*domain.Payment) *[]CreatePaymentOutput {
	outputs := make([]CreatePaymentOutput, len(payments))
	for i, payment := range payments {
		outputs[i] = *FromPayment(payment)
	}
	return &outputs
}
