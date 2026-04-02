package usecase

import (
	"context"
	"fmt"
	"payment-processor/internal/payment/adapters/inbound/dto"
	"payment-processor/internal/payment/ports"
)

type PaymentGetUseCaseInterface interface {
	Execute(ctx context.Context) (*[]dto.CreatePaymentOutput, error)
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

func (uc *PaymentGetUseCase) Execute(ctx context.Context) (*[]dto.CreatePaymentOutput, error) {

	payments, err := uc.paymentRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payments: %w", err)
	}

	return dto.FromPayments(payments), nil
}
