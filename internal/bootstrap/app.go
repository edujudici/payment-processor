package bootstrap

import (
	"log"
	"log/slog"
	"os"
	"payment-processor/internal/payment/adapters/inbound/handler"
	"payment-processor/internal/payment/adapters/outbound/repository"
	"payment-processor/internal/payment/adapters/outbound/service"
	"payment-processor/internal/payment/usecase"
)

type App struct {
	Handlers *Handlers
}

type Handlers struct {
	Payment    *handler.PaymentProcessorHandler
	Preference *handler.PreferenceHandler
}

func NewApp(logger *slog.Logger) *App {
	// =========================
	// ENV
	// =========================
	mpToken := os.Getenv("MERCADO_PAGO_TOKEN")
	if mpToken == "" {
		log.Fatal("MERCADO_PAGO_TOKEN not set")
	}

	// =========================
	// Infra
	// =========================
	db := repository.NewMySQLConnection()

	// =========================
	// Payment - Adapters outbound
	// =========================
	paymentRepo := repository.NewPaymentProcessorRepositoryMySQL(db)
	paymentSvc := service.NewPaymentCancelService(nil)

	// =========================
	// Payment - Usecases
	// =========================
	payCreateUC := usecase.NewPaymentCreateUseCase(paymentSvc, paymentRepo)
	payGetUC := usecase.NewPaymentGetUseCase(paymentSvc, paymentRepo)

	// =========================
	// Payment - Handler
	// =========================
	paymentHandler := handler.NewPaymentProcessorHandler(payCreateUC, payGetUC)

	// =========================
	// Preference - Adapters outbound
	// =========================
	preferenceRepo := repository.NewPreferenceRepositoryMySQL(db)
	preferenceSvc, err := service.NewPreferenceCreateService(mpToken)
	if err != nil {
		log.Fatalf("Failed to create PreferenceCreateService: %v", err)
	}

	// =========================
	// Preference - Usecases
	// =========================
	preferenceCreateUC := usecase.NewPreferenceCreateUseCase(preferenceSvc, preferenceRepo)

	// =========================
	// Preference - Handler
	// =========================
	preferenceHandler := handler.NewPreferenceHandler(preferenceCreateUC)

	return &App{
		Handlers: &Handlers{
			Payment:    paymentHandler,
			Preference: preferenceHandler,
		},
	}
}
