package logging

import (
	"log/slog"
	"os"

	"github.com/golang-auth/internal/core/ports"
)

type stdoutAdapter struct {
	logger *slog.Logger
}

func NewStdoutLogger() ports.Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	return &stdoutAdapter{
		logger: slog.New(handler),
	}
}

func (l *stdoutAdapter) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *stdoutAdapter) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *stdoutAdapter) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *stdoutAdapter) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *stdoutAdapter) Fatal(msg string, args ...any) {
	l.logger.Error("Fatal "+msg, args...)
	os.Exit(1)
}

// // WithContext is the most critical part for Hexagonal/Tracing.
// // It creates a new multi-logger where every sub-logger is trace-aware.
// func (m *multi) WithContext(ctx context.Context) ports.Logger {
// 	contextualLoggers := make([]ports.Logger, len(m.loggers))
// 	for i, logger := range m.loggers {
// 		contextualLoggers[i] = logger.WithContext(ctx)
// 	}
// 	return &multi{loggers: contextualLoggers}
// }
