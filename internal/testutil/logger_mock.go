package testutil

type NoopLogger struct{}

func (l *NoopLogger) Debug(msg string, args ...interface{}) {}
func (l *NoopLogger) Info(msg string, args ...interface{})  {}
func (l *NoopLogger) Error(msg string, args ...interface{}) {}
func (l *NoopLogger) Warn(msg string, args ...interface{})  {}
func (l *NoopLogger) Fatal(msg string, args ...interface{}) {}
