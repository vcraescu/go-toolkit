package log

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}

var _ Logger = (*JSONLogger)(nil)

type JSONLogger struct {
	logger *slog.Logger
	ctx    context.Context
}

func New() *JSONLogger {
	return &JSONLogger{
		logger: slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
		})),
		ctx: context.Background(),
	}
}

func (l JSONLogger) With(args ...any) Logger {
	newArgs := make([]any, 0, len(args))

	for _, arg := range args {
		if ctx, ok := arg.(context.Context); ok {
			l.ctx = ctx

			continue
		}

		newArgs = append(newArgs, arg)
	}

	return &JSONLogger{
		logger: l.logger.With(newArgs...),
		ctx:    l.ctx,
	}
}

func (l JSONLogger) Info(msg string, args ...any) {
	l.logger.InfoContext(l.ctx, msg, l.withTrace(l.ctx, args)...)
}

func (l JSONLogger) Error(msg string, args ...any) {
	l.logger.ErrorContext(l.ctx, msg, l.withTrace(l.ctx, args)...)
}

func (l JSONLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, l.withTrace(ctx, args)...)
}

func (l JSONLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, l.withTrace(ctx, args)...)
}

func (l JSONLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, l.withTrace(l.ctx, args)...)
}

func (l JSONLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.Warn(msg, l.withTrace(ctx, args)...)
}

func (l JSONLogger) withTrace(ctx context.Context, args []any) []any {
	spanCtx := trace.SpanContextFromContext(ctx)

	if spanCtx.TraceID().IsValid() {
		args = append(args, slog.String("traceID", spanCtx.TraceID().String()))
	}

	if spanCtx.SpanID().IsValid() {
		args = append(args, slog.String("spanID", spanCtx.SpanID().String()))
	}

	return args
}

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}
