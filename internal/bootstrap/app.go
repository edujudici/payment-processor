package bootstrap

import (
	"log/slog"
	"payment-processor/internal/payment/adapters/inbound/handler"
	"payment-processor/internal/payment/adapters/outbound/repository"
	"payment-processor/internal/payment/adapters/outbound/service"
	"payment-processor/internal/payment/usecase"
)

type App struct {
	PaymentHandler *handler.PaymentProcessorHandler
}

func NewApp(logger *slog.Logger) *App {
	// Infra
	db := repository.NewMySQLConnection()

	// Adapters outbound
	paymentRepo := repository.NewPaymentProcessorRepositoryMySQL(db)
	paymentSvc := service.NewPaymentCancelService(nil)

	// Usecase
	uc := usecase.NewPaymentUseCase(paymentSvc, paymentRepo)

	// Handler (inbound)
	h := handler.NewPaymentProcessorHandler(uc)

	return &App{
		PaymentHandler: h,
	}
}
