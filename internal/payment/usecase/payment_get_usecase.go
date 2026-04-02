package usecase

import (
	"context"
	"fmt"
	"payment-processor/internal/payment/ports"
)

type PaymentGetUseCaseInterface interface {
	Execute(ctx context.Context) (any, error)
}

type PaymentGetUseCase struct {
	paymentService    ports.PaymentSvcInterface
	paymentRepository ports.PaymentRepository
}

func NewPaymentGetUseCase(
	paymentService ports.PaymentSvcInterface,
	paymentRepository ports.PaymentRepository,
) *PaymentGetUseCase {
	return &PaymentGetUseCase{
		paymentService:    paymentService,
		paymentRepository: paymentRepository,
	}
}

func (uc *PaymentGetUseCase) Execute(ctx context.Context) (any, error) {

	payments, err := uc.paymentRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payments: %w", err)
	}

	return payments, nil
}
