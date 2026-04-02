package usecase

import (
	"context"
	"fmt"
	"payment-processor/internal/payment/adapters/inbound/dto"
	"payment-processor/internal/payment/domain"
	"payment-processor/internal/payment/ports"
	"time"
)

type PaymentCreateUseCaseInterface interface {
	Execute(ctx context.Context, req dto.Request) error
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

func (uc *PaymentCreateUseCase) Execute(ctx context.Context, req dto.Request) error {

	payment := &domain.Payment{
		ID:        "generated-id",
		Amount:    req.Amount,
		Status:    "created",
		CreatedAt: time.Now(),
	}

	err := uc.paymentRepository.Save(ctx, payment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}
