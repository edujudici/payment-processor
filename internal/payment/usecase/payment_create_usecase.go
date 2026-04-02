package usecase

import (
	"context"
	"fmt"
	"payment-processor/internal/payment/adapters/inbound/dto"
	"payment-processor/internal/payment/ports"
)

type PaymentCreateUseCaseInterface interface {
	Execute(ctx context.Context, input dto.CreatePaymentInput) (*dto.CreatePaymentOutput, error)
}

type PaymentCreateUseCase struct {
	paymentService    ports.PaymentSvcInterface
	paymentRepository ports.PaymentRepository
}

func NewPaymentCreateUseCase(
	paymentService ports.PaymentSvcInterface,
	paymentRepository ports.PaymentRepository,
) *PaymentCreateUseCase {
	return &PaymentCreateUseCase{
		paymentService:    paymentService,
		paymentRepository: paymentRepository,
	}
}

func (uc *PaymentCreateUseCase) Execute(ctx context.Context, input dto.CreatePaymentInput) (*dto.CreatePaymentOutput, error) {

	payment, err := dto.ToPayment(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	err = uc.paymentRepository.Save(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return dto.FromPayment(payment), nil
}
