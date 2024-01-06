package echo

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/vcraescu/go-toolkit/log"
	"github.com/vcraescu/go-toolkit/transport/echo/middleware"
	"net/http"
)

type Server struct {
	*echo.Echo

	logger log.Logger
}

func NewServer(opts ...ServerOption) *Server {
	options := newServerOptions(opts...)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	middlewares := make([]echo.MiddlewareFunc, 0)

	if options.tracer != nil {
		middlewares = append(middlewares, middleware.WithTracer(options.tracer))
	}

	if options.logger != nil {
		middlewares = append(middlewares, middleware.WithLogger(options.logger))
	}

	e.Use(middlewares...)

	if options.healthHandler != nil {
		e.GET("/health", options.healthHandler)
	}

	return &Server{
		Echo:   e,
		logger: options.logger,
	}
}

func (s *Server) Start(ctx context.Context, address string) error {
	go s.handleShutdown(ctx)

	s.logger.Info(ctx, "server started", log.String("address", address))

	if err := s.Echo.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) handleShutdown(ctx context.Context) {
	<-ctx.Done()

	s.logger.Info(ctx, "server is shutting down...")

	if err := s.Shutdown(ctx); err != nil {
		s.logger.Error(ctx, "server shutdown failed", log.Error(err))
	}
}
