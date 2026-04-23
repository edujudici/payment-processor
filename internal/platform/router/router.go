package router

import (
	"net/http"
	"payment-processor/internal/bootstrap"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(app *bootstrap.App) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Routes
	// r.Route("/api/v1", func(r chi.Router) {
	r.Get("/payments", app.Handlers.Payment.GetPayments)
	r.Post("/payments", app.Handlers.Payment.CreatePayment)

	r.Post("/preference", app.Handlers.Preference.CreatePreference)
	// })

	return r
}
