package ports

import (
	"context"
	"payment-processor/internal/payment/domain"
)

type PaymentRepository interface {
	FindAll(ctx context.Context) ([]*domain.Payment, error)
	FindByID(ctx context.Context, id string) (*domain.Payment, error)
	Save(ctx context.Context, payment *domain.Payment) error
}
