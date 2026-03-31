package ports

import (
	"boleto-cancel/internal/boleto/domain"
	"context"
)

type BoletoCancelRepository interface {
	GetBank(ctx context.Context, bankId string) (*domain.Bank, error)
	GetBusiness(ctx context.Context, businessId string) (*domain.Business, error)
	GetBankByAgreement(ctx context.Context, bankAgreement string) (*domain.Bank, error)
}
