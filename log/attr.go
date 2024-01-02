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

type Attr interface {
	Key() string
	Value() any
	getInternal() slog.Attr
}

var _ Attr = (*attr)(nil)

type attr struct {
	internal slog.Attr
}

func newAttr(a slog.Attr) Attr {
	return attr{internal: a}
}

func toSlogAttrs(attrs ...any) []any {
	for i, attr := range attrs {
		if a, ok := attr.(Attr); ok {
			attrs[i] = a.getInternal()
		}
	}

	return attrs
}

func (a attr) getInternal() slog.Attr {
	return a.internal
}

func (a attr) Key() string {
	return a.internal.Key
}

func (a attr) Value() any {
	return a.internal.Value
}

func Error(err error) Attr {
	return newAttr(slog.String(ErrorKey, err.Error()))
}

func traceAttrs(ctx context.Context) []slog.Attr {
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

func Any(key string, value any) Attr {
	return newAttr(slog.Any(key, value))
}

func Int(key string, value int) Attr {
	return newAttr(slog.Int(key, value))
}

func Int64(key string, value int64) Attr {
	return newAttr(slog.Int64(key, value))
}

func Float64(key string, value float64) Attr {
	return newAttr(slog.Float64(key, value))
}

func Float(key string, value float32) Attr {
	return newAttr(slog.Float64(key, float64(value)))
}

func Bool(key string, value bool) Attr {
	return newAttr(slog.Bool(key, value))
}

func String(key string, value string) Attr {
	return newAttr(slog.String(key, value))
}
