package httpserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-auth/internal/adapters/config"
	"github.com/golang-auth/internal/adapters/logging"
	"github.com/golang-auth/internal/adapters/repository/postgre"
	"github.com/golang-auth/internal/core/ports"
	"github.com/redis/go-redis/v9"
)

func Run(ctx context.Context, logger ports.Logger, handler http.Handler, addr string, serverName string) error {
	s := &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		logger.Info("Starting HTTP "+serverName+" server", "address", s.Addr)
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP "+serverName+" server failed", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down HTTP " + serverName + " server...")

	// Give the server 5 seconds to finish processing existing requests
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.Shutdown(shutdownCtx)
}

func LoadComponents() (ports.Logger, *postgre.Client, *redis.Client) {
	// Configuration
	cfg, err := config.NewLoggingConfig()
	if err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	// Logger
	logger := logging.NewLogger(cfg)
	logger.Info("Logging successfully configured to use the adapter: ", cfg.Adapter())

	// PostgreSQL
	logger.Info("Loading PostgreSQL config")
	postgreConfig, err := config.NewDefaultDBConfig()
	if err != nil {
		logger.Error("Failed to load PostgreSQL config", "error", err)
		os.Exit(1)
	}

	logger.Info("Connecting to PostgreSQL database")
	client, err := postgre.NewPostgreSQLClient(postgreConfig)
	if err != nil {
		logger.Error("Postgresql connection error", "error", err)
		os.Exit(1)
	}
	logger.Info("Successful PostgreSQL connection")

	// Redis
	logger.Info("Connecting to redis server")
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	logger.Info("Successful redis connection")

	return logger, client, rdb
}
