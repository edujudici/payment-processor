package usecase

import (
	"context"
	"fmt"
	"payment-processor/internal/payment/adapters/inbound/dto"
	"payment-processor/internal/payment/domain"
	"payment-processor/internal/payment/ports"
	"time"
)

type PaymentUseCaseInterface interface {
	GetPayments(ctx context.Context) (any, error)
	CreatePayment(ctx context.Context, req dto.Request) error
}

type PaymentUseCase struct {
	paymentService    ports.PaymentSvcInterface
	paymentRepository ports.PaymentRepository
}

func NewPaymentUseCase(
	paymentService ports.PaymentSvcInterface,
	paymentRepository ports.PaymentRepository,
) *PaymentUseCase {
	return &PaymentUseCase{
		paymentService:    paymentService,
		paymentRepository: paymentRepository,
	}
}

func (uc *PaymentUseCase) GetPayments(ctx context.Context) (any, error) {

	payments, err := uc.paymentRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payments: %w", err)
	}

	return payments, nil
}

func (uc *PaymentUseCase) CreatePayment(ctx context.Context, req dto.Request) error {

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
