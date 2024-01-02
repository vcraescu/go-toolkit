package log

import "context"

var _ Logger = (*NopLogger)(nil)

type NopLogger struct{}

func (l NopLogger) With(args ...any) Logger {
	return l
}

func NewNopLogger() *NopLogger {
	return &NopLogger{}
}

func (l NopLogger) Info(ctx context.Context, msg string, args ...any) {}

func (l NopLogger) Error(ctx context.Context, msg string, args ...any) {}

func (l NopLogger) Warn(ctx context.Context, msg string, args ...any) {}

func (l NopLogger) Debug(ctx context.Context, msg string, args ...any) {}
