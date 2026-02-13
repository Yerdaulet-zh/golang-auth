package logging

// import (
// 	"context"
// 	"log/slog"
// 	"os"

// 	"github.com/go-hexagonal-practice/internal/core/ports"
// 	"go.opentelemetry.io/contrib/bridges/otelslog"
// 	"go.opentelemetry.io/otel/log/global"
// )

// type otelAdapter struct {
// 	logger *slog.Logger
// 	ctx    context.Context // Store context to allow slog to extract TraceIDs
// }

// func NewOtelLogger(serviceName string) ports.Logger {
// 	return &otelAdapter{
// 		logger: otelslog.NewLogger(serviceName),
// 		ctx:    context.Background(),
// 	}
// }

// func (a *otelAdapter) Debug(msg string, args ...any) {
// 	a.logger.DebugContext(a.ctx, msg, args...)
// }

// func (a *otelAdapter) Info(msg string, args ...any) {
// 	a.logger.InfoContext(a.ctx, msg, args...)
// }

// func (a *otelAdapter) Warn(msg string, args ...any) {
// 	a.logger.WarnContext(a.ctx, msg, args...)
// }

// func (a *otelAdapter) Error(msg string, args ...any) {
// 	a.logger.ErrorContext(a.ctx, msg, args...) // Fixed method call
// }

// func (a *otelAdapter) Fatal(msg string, args ...any) {
// 	a.logger.ErrorContext(a.ctx, "FATAL: "+msg, args...) // Log as error before exit
// 	os.Exit(1)
// }

// func (a *otelAdapter) WithContext(ctx context.Context) ports.Logger {
// 	// We return a new adapter instance that holds the specific request context
// 	return &otelAdapter{
// 		logger: otelslog.NewLogger("dice", otelslog.WithLoggerProvider(global.GetLoggerProvider())),
// 		ctx:    ctx,
// 	}
// }
