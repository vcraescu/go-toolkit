package log

import (
	"context"
	"log/slog"
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}

type JSONLogger struct {
	*slog.Logger
}

func New(opts ...Option) *JSONLogger {
	options := newOptions(opts...)

	h := NewHandler(slog.NewJSONHandler(options.output, &slog.HandlerOptions{
		Level: slog.Level(options.level),
	}))
	h.setClock(options.clock)

	return &JSONLogger{
		Logger: slog.New(h),
	}
}

func (l *JSONLogger) With(args ...any) Logger {
	return &JSONLogger{
		Logger: l.Logger.With(args...),
	}
}

func (l *JSONLogger) Info(ctx context.Context, msg string, args ...any) {
	l.convertAttrs(args)

	l.Logger.InfoContext(ctx, msg, args...)
}

func (l *JSONLogger) Error(ctx context.Context, msg string, args ...any) {
	l.convertAttrs(args)

	l.Logger.ErrorContext(ctx, msg, args...)
}

func (l *JSONLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.convertAttrs(args)

	l.Logger.WarnContext(ctx, msg, args...)
}

func (l *JSONLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.convertAttrs(args)

	l.Logger.DebugContext(ctx, msg, args...)
}

func (l *JSONLogger) convertAttrs(args []any) {
	for i, arg := range args {
		if a, ok := arg.(Attr); ok {
			args[i] = a.getInternal()
		}
	}
}
