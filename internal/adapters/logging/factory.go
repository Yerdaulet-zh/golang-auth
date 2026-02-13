package logging

import (
	"github.com/golang-auth/internal/adapters/config"
	"github.com/golang-auth/internal/core/ports"
)

func NewLogger(cfg *config.LoggingConfig) ports.Logger {
	switch cfg.Adapter() {
	case "loki":
		return NewLokiLogger(
			cfg.LokiURL(),
			cfg.LokiLabels(),
		)
	case "multi":
		stdoutLogger := NewStdoutLogger()
		lokiLogger := NewLokiLogger(
			cfg.LokiURL(),
			cfg.LokiLabels(),
		)
		loggers := []ports.Logger{stdoutLogger, lokiLogger}
		return NewMultiLogger(loggers...)
	default:
		return NewStdoutLogger()
	}
}
