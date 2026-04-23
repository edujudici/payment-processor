package ports

import (
	"context"
	"payment-processor/internal/payment/domain"
)

type PreferenceRepository interface {
	Save(ctx context.Context, preference *domain.Preference) error
}
