package middleware

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vcraescu/go-toolkit/log"
	"go.opentelemetry.io/otel/trace"
)

func WithLogger(logger log.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			ctx := c.Request().Context()
			logger := logger.With(
				log.String("uri", v.URI),
				log.String("method", v.Method),
				log.Int("status", v.Status),
			)

			if v.Error != nil {
				logger.Error(ctx, "REQ_ERR", log.Error(v.Error))

				return nil
			}

			logger.Info(ctx, "REQ")

			return nil
		},
	})
}

func WithTracer(tracer trace.Tracer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, span := tracer.Start(c.Request().Context(), c.Request().RequestURI)
			defer span.End()

			c.SetRequest(c.Request().WithContext(context.WithoutCancel(ctx)))

			return next(c)
		}
	}
}
