package ports

import (
	"context"
	"payment-processor/internal/payment/adapters/outbound/service"
)

type PreferenceSvcInterface interface {
	PreferenceCreate(ctx context.Context, req *service.PreferenceCreateInput) (*service.PreferenceCreateOutput, error)
}
