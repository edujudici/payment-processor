package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"payment-processor/internal/bootstrap"
	"payment-processor/internal/platform/router"
	"payment-processor/internal/platform/server"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// App (DI)
	app := bootstrap.NewApp(logger)

	// Router
	r := router.NewRouter(app)

	// Server
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := os.Getenv("PORT")
	srv := server.NewHTTPServer(r, port, logger)

	// Start
	if err := srv.Start(ctx); err != nil {
		logger.Error("server failed", "error", err)
	}
}
