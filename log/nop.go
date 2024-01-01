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

func (l NopLogger) InfoContext(ctx context.Context, msg string, args ...any) {}

func (l NopLogger) ErrorContext(ctx context.Context, msg string, args ...any) {}

func (l NopLogger) Info(msg string, args ...any) {}

func (l NopLogger) Error(msg string, args ...any) {}

func (l NopLogger) Warn(msg string, args ...any) {}

func (l NopLogger) WarnContext(ctx context.Context, msg string, args ...any) {}
