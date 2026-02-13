package httpserver

import (
	"net/http"
	"time"

	http_hanlder "github.com/golang-auth/internal/adapters/handlers/http"
	"github.com/golang-auth/internal/adapters/handlers/http/middleware"
	"github.com/golang-auth/internal/adapters/repository/postgre"
	"github.com/golang-auth/internal/core/ports"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

type Middleware func(http.Handler) http.Handler

// Create a helper to chain them
func ApplyMiddleware(h http.Handler, mws ...Middleware) http.Handler {
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}

func MapBusinessRoutes(
	logger ports.Logger,
	rdb *redis.Client,
	userService ports.UserUseCase,
) http.Handler {
	mux := http.NewServeMux()

	userHandler := http_hanlder.NewUserHandler(userService, logger)
	mux.HandleFunc("POST /v1/register", userHandler.Register)
	// mux.HandleFunc("POST /v1/register/verify/{token}", userHandler.VerifyUserEmail)
	// mux.HandleFunc("POST /v1/login", userHandler.Login)

	middlewares := []Middleware{
		middleware.LoggingMiddleware(logger),                   // 3. Log everything (including blocks)
		middleware.IPRateLimiter(logger, rdb, 10, time.Minute), // 2. Then check limit
		// middleware.RecoveryMiddleware(logger),               // 1. Catch panics first
	}
	return ApplyMiddleware(mux, middlewares...)
}

func MapManagementRoutes(logger ports.Logger, db *postgre.Client, reg *prometheus.Registry) http.Handler {
	mux := http.NewServeMux()

	healthHdl := NewHealthHandler(db)
	mux.HandleFunc("GET /healthz", healthHdl.Healthz)
	mux.HandleFunc("GET /ready", healthHdl.Ready)

	mux.Handle("GET /metrics", promhttp.Handler())
	return mux
}
