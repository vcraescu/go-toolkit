package log

import (
	"context"
	"fmt"
	"log/slog"
	"path"
	"runtime"
	"time"
)

var _ slog.Handler = (*DefaultHandler)(nil)

type Handler interface {
	slog.Handler
}

type DefaultHandler struct {
	handler Handler
	clock   time.Time
}

func NewHandler(handler Handler) *DefaultHandler {
	return &DefaultHandler{
		handler: handler,
	}
}

func (h *DefaultHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *DefaultHandler) Handle(ctx context.Context, record slog.Record) error {
	if !h.clock.IsZero() {
		record.Time = h.clock
	}

	attrs := append(traceAttrs(ctx), h.getSourceAttr(record.PC))
	record.AddAttrs(attrs...)

	return h.handler.Handle(ctx, record)
}

func (h *DefaultHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &DefaultHandler{handler: h.handler.WithAttrs(attrs), clock: h.clock}
}

func (h *DefaultHandler) WithGroup(name string) slog.Handler {
	return &DefaultHandler{handler: h.handler.WithGroup(name), clock: h.clock}
}

func (h *DefaultHandler) getSourceAttr(pc uintptr) slog.Attr {
	rpc := []uintptr{pc}

	runtime.Callers(6, rpc)

	fs := runtime.CallersFrames(rpc)
	f, _ := fs.Next()

	if f.File == "" {
		return slog.Attr{}
	}

	file := path.Clean(path.Join(path.Base(path.Dir(f.File)), path.Base(f.File)))

	return slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", file, f.Line))
}

func (h *DefaultHandler) setClock(clock time.Time) {
	h.clock = clock
}
