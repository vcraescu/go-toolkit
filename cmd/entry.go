package cmd

import (
	"context"
	"fmt"
	"github.com/vcraescu/go-toolkit/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"os"
	"os/signal"
	"time"
)

var termCh chan os.Signal

type shutdowner interface {
	Shutdown(ctx context.Context) error
}

type Handler interface {
	Handle(ctx Context) error
}

var _ Handler = HandlerFunc(nil)

type HandlerFunc func(ctx Context) error

func (fn HandlerFunc) Handle(ctx Context) error {
	return fn(ctx)
}

func Start(handler Handler, opts ...Option) {
	options := newOptions(opts...)

	otel.SetTracerProvider(options.traceProvider)

	ctx := &startContext{
		context: options.ctx,
		logger:  options.logger,
		tracer:  otel.Tracer("app"),
		config:  options.config,
	}

	if v, ok := options.traceProvider.(shutdowner); ok {
		defer v.Shutdown(ctx.context)
	}

	if options.appName != "" {
		ctx.logger = options.logger.With(log.String("app", options.appName))
	}

	if options.terminateTimeout > 0 {
		ctx.context = TerminateContext(ctx.context, options.terminateTimeout, ctx.logger)
	}

	if options.appName != "" {
		ctx.tracer = otel.Tracer(options.appName)
	}

	var span trace.Span

	ctx.context, span = ctx.tracer.Start(ctx.context, "run")
	defer span.End()

	if err := handler.Handle(ctx); err != nil {
		panic(fmt.Errorf("handle: %w", err))
	}
}

func TerminateContext(ctx context.Context, timeout time.Duration, logger log.Logger) context.Context {
	ch := termCh

	if ch == nil {
		ch = make(chan os.Signal, 1)
	}

	ctx, cancel := context.WithCancel(ctx)
	signal.Notify(ch, os.Interrupt)

	go func() {
		defer cancel()
		sig := <-ch

		logger.Info(
			ctx,
			"terminating...",
			log.String("sig", sig.String()),
			log.String("timeout", timeout.String()),
		)

		time.Sleep(timeout)
	}()

	return ctx
}
