package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"payment-processor/internal/payment/adapters/inbound/dto"
	"payment-processor/internal/payment/adapters/outbound/service"
	"payment-processor/internal/payment/ports"
	"time"
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

	url := os.Getenv("URL")
	if url == "" {
		log.Fatal("URL not set")
	}

	protocol := generateExternalReference("protocol")

	pref, err := uc.paymentService.PaymentCreate(ctx, &service.PaymentCreateInput{
		Item:  dto.MapItem(input.Item),
		Payer: dto.MapPayer(input.Payer),
		BackURLs: service.BackURLs{
			Success: input.BackURL + "/success",
			Failure: input.BackURL + "/failure",
			Pending: input.BackURL + "/pending",
		},
		AutoReturn:        "approved",
		NotificationURL:   url + "/notification",
		ExternalReference: protocol,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate new payment: %w", err)
	}

	payment, err := dto.ToPayment(
		protocol,
		input.Payer.Name,
		input.Payer.Surname,
		input.Payer.Email,
		input.Item.Quantity,
		input.Item.UnitPrice*float64(input.Item.Quantity),
		input.Item.UnitPrice*float64(input.Item.Quantity),
		input.Item.Title,
		pref.ID,
		pref.ExternalReference,
		pref.InitPoint,
		pref.SandboxInitPoint,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to convert dto to domain payment: %w", err)
	}

	err = uc.paymentRepository.Save(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	return dto.FromPayment(payment), nil
}

func generateExternalReference(prefix string) string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)

	return fmt.Sprintf("%s-%d-%s", prefix, time.Now().Unix(), hex.EncodeToString(b))
}
