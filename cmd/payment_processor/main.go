package main

import (
	"log"
	"log/slog"
	"os"
	"payment-processor/internal/bootstrap"
	"payment-processor/internal/platform/router"
	"payment-processor/internal/platform/server"

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
	port := os.Getenv("PORT")
	srv := server.NewHTTPServer(r, port, logger)

	// Start
	srv.Start()
}
