package cmd

import (
	"context"
	"github.com/vcraescu/go-toolkit/config"
	"github.com/vcraescu/go-toolkit/log"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"time"
)

type options struct {
	terminateTimeout time.Duration
	appName          string
	config           config.Getter
	traceProvider    trace.TracerProvider
	logger           log.Logger
	ctx              context.Context
}

type Option interface {
	apply(options *options)
}

var _ Option = optionFunc(nil)

type optionFunc func(options *options)

func (fn optionFunc) apply(options *options) {
	fn(options)
}

func WithTerminateTimeout(timeout time.Duration) Option {
	return optionFunc(func(options *options) {
		options.terminateTimeout = timeout
	})
}

func WithAppName(appName string) Option {
	return optionFunc(func(options *options) {
		options.appName = appName
	})
}

func WithConfig(cfg config.Getter) Option {
	return optionFunc(func(options *options) {
		options.config = cfg
	})
}

func WithLogger(logger log.Logger) Option {
	return optionFunc(func(options *options) {
		options.logger = logger
	})
}

func WithTraceProvider(traceProvider trace.TracerProvider) Option {
	return optionFunc(func(options *options) {
		options.traceProvider = traceProvider
	})
}

func WithContext(ctx context.Context) Option {
	return optionFunc(func(options *options) {
		options.ctx = ctx
	})
}

func newOptions(opts ...Option) *options {
	options := &options{
		logger:        log.New(),
		traceProvider: noop.NewTracerProvider(),
		ctx:           context.Background(),
	}

	for _, opt := range opts {
		opt.apply(options)
	}

	if options.logger == nil {
		options.logger = log.NewNop()
	}

	return options
}
