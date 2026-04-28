package ports

import (
	"context"
	"payment-processor/internal/payment/adapters/outbound/service"
)

type PaymentSvcInterface interface {
	PaymentCreate(ctx context.Context, req *service.PaymentCreateInput) (*service.PaymentCreateOutput, error)
}
