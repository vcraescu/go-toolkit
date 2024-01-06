package echo

import (
	"github.com/labstack/echo/v4"
	"github.com/vcraescu/go-toolkit/log"
	"github.com/vcraescu/go-toolkit/transport/echo/handler"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type serverOptions struct {
	tracer        trace.Tracer
	logger        log.Logger
	healthHandler echo.HandlerFunc
}

type ServerOption interface {
	apply(opts *serverOptions)
}

var _ ServerOption = serverOptionFunc(nil)

type serverOptionFunc func(opts *serverOptions)

func (fn serverOptionFunc) apply(opts *serverOptions) {
	fn(opts)
}

func WithTracer(tracer trace.Tracer) ServerOption {
	return serverOptionFunc(func(opts *serverOptions) {
		opts.tracer = tracer
	})
}

func WithHealth(h echo.HandlerFunc) ServerOption {
	return serverOptionFunc(func(opts *serverOptions) {
		opts.healthHandler = h
	})
}

func newServerOptions(opts ...ServerOption) *serverOptions {
	options := &serverOptions{
		tracer:        noop.NewTracerProvider().Tracer("noop"),
		logger:        log.New(),
		healthHandler: handler.HealthHandler,
	}

	for _, opt := range opts {
		opt.apply(options)
	}

	return options
}
