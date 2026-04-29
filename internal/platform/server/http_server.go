package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type HTTPServer struct {
	server *http.Server
	logger *slog.Logger
}

func NewHTTPServer(handler http.Handler, port string, logger *slog.Logger) *HTTPServer {
	if port == "" {
		port = "8080"
	}

	handler = withCORS(handler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &HTTPServer{
		logger: logger,
		server: server,
	}
}

func withCORS(next http.Handler) http.Handler {
	allowedOrigins := map[string]bool{
		"https://buscacep.escaliagora.com.br": true,
		"http://localhost:3000":               true,
		"http://localhost:5173":               true, // Vite
		"http://127.0.0.1:3000":               true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *HTTPServer) Start(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		s.logger.Info("server starting", "addr", s.server.Addr)

		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("shutdown signal received")
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.logger.Info("shutting down server...")

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	s.logger.Info("server stopped")
	return nil
}
