package logging

import (
	"os"

	"github.com/golang-auth/internal/core/ports"
)

type multi struct {
	loggers []ports.Logger
}

func NewMultiLogger(loggers ...ports.Logger) ports.Logger {
	return &multi{
		loggers: loggers,
	}
}

func (m *multi) Debug(msg string, args ...any) {
	for _, logger := range m.loggers {
		logger.Debug(msg, args...)
	}
}

func (m *multi) Info(msg string, args ...any) {
	for _, logger := range m.loggers {
		logger.Info(msg, args...)
	}
}

func (m *multi) Warn(msg string, args ...any) {
	for _, logger := range m.loggers {
		logger.Debug(msg, args...)
	}
}

func (m *multi) Error(msg string, args ...any) {
	for _, logger := range m.loggers {
		logger.Error(msg, args...)
	}
}

func (m *multi) Fatal(msg string, args ...any) {
	for _, logger := range m.loggers {
		logger.Fatal(msg, args...)
		os.Exit(1)
	}
}

// // WithContext allows slog to use context-based attributes if you add them later
// func (l *stdoutAdapter) WithContext(ctx context.Context) ports.Logger {
// 	return &stdoutAdapter{
// 		logger: l.logger,
// 	}
// }
