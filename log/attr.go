package log

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

const (
	ErrorKey   = "error"
	TraceIDKey = "traceID"
	SpanIDKey  = "spanID"
)

func Error(err error) slog.Attr {
	return slog.String(ErrorKey, err.Error())
}

func Trace(ctx context.Context) []slog.Attr {
	spanCtx := trace.SpanContextFromContext(ctx)

	attrs := make([]slog.Attr, 0)

	if spanCtx.TraceID().IsValid() {
		attrs = append(attrs, slog.String(TraceIDKey, spanCtx.TraceID().String()))
	}

	if spanCtx.SpanID().IsValid() {
		attrs = append(attrs, slog.String(SpanIDKey, spanCtx.SpanID().String()))
	}

	return attrs
}

func Any(key string, value any) slog.Attr {
	return slog.Any(key, value)
}
