package main

import (
	"context"
	"os/signal"
	"syscall"

	httpserver "github.com/golang-auth/cmd/http_server"
	"github.com/golang-auth/internal/adapters/config"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger, client, rdb := httpserver.LoadComponents()

	defer func() {
		logger.Info("Closing infrastructure connections...")
		if err := client.Close(); err != nil {
			logger.Error("Postgre close error", "error", err)
		}
		if err := rdb.Close(); err != nil {
			logger.Error("Redis close error", "error", err)
		}
		logger.Info("Done")
	}()

	logger.Info("Loading HTTP Server config")
	httpConfig := config.NewHttpConfig()
	logger.Info("Successfully loaded HTTP Server config")

	reg := prometheus.NewRegistry()

	// userRepo := postgre.NewUserRepository(client.DB, logger)
	// userService := service.NewUserService(userRepo)

	// mapBusinessHandler := httpserver.MapBusinessRoutes(logger, rdb, userService)
	mapManagementRoutes := httpserver.MapManagementRoutes(logger, client, reg)
	errChan := make(chan error, 1)

	// go func() {
	// 	errChan <- httpserver.Run(ctx, logger, mapBusinessHandler, ...)
	// }()

	go func() {
		errChan <- httpserver.Run(ctx, logger, mapManagementRoutes, httpConfig.HttpManagementAddr(), "Management")
	}()

	select {
	case err := <-errChan:
		if err != nil {
			logger.Error("Critical server failure", "error", err)
			// This triggers the shutdown of all other goroutines via context
			stop()
		}
	case <-ctx.Done():
		logger.Info("Shutdown requested by user")
	}

	logger.Info("Application exited cleanly")
}
