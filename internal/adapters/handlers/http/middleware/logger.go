package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-auth/internal/core/ports"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(logger ports.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{w, http.StatusOK}
			next.ServeHTTP(rw, r)

			// Log the details
			logger.Info("request_processed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rw.statusCode),
				slog.Duration("latency", time.Since(start)),
				slog.String("ip", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)
		})
	}
}
