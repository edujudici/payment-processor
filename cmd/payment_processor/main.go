package main

import (
	"log/slog"
	"os"
	"payment-processor/internal/bootstrap"
	"payment-processor/internal/platform/router"
	"payment-processor/internal/platform/server"
)

func main() {
	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// App (DI)
	app := bootstrap.NewApp(logger)

	// Router
	r := router.NewRouter(app)

	// Server
	port := os.Getenv("PORT")
	srv := server.NewHTTPServer(r, port, logger)

	// Start
	srv.Start()
}
