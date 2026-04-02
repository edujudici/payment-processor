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

	// Usecases
	payCreateUC := usecase.NewPaymentCreateUseCase(paymentSvc, paymentRepo)
	payGetUC := usecase.NewPaymentGetUseCase(paymentSvc, paymentRepo)

	// Handler (inbound)
	h := handler.NewPaymentProcessorHandler(payCreateUC, payGetUC)

	return &App{
		PaymentHandler: h,
	}
}
