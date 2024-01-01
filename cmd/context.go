package cmd

import (
	"context"
	"github.com/vcraescu/go-toolkit/config"
	"github.com/vcraescu/go-toolkit/log"
	"go.opentelemetry.io/otel/trace"
)

type Context interface {
	Logger() log.Logger
	Tracer() trace.Tracer
	Config() config.Getter
	Context() context.Context
}

var _ Context = (*startContext)(nil)

type startContext struct {
	logger  log.Logger
	tracer  trace.Tracer
	config  config.Getter
	context context.Context
}

func (c *startContext) Context() context.Context {
	return c.context
}

func (c *startContext) Config() config.Getter {
	return c.config
}

func (c *startContext) Logger() log.Logger {
	return c.logger
}

func (c *startContext) Tracer() trace.Tracer {
	return c.tracer
}
