package log_test

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/go-toolkit/log"
	"github.com/vcraescu/go-toolkit/testing/jsontest"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := trace.ContextWithSpanContext(context.Background(), newSpanContext())
		buf := &bytes.Buffer{}

		logger := log.New(log.WithOutput(buf), log.WithClock(now)).
			With(log.String("a", "b"))

		logger.Info(ctx, "test", log.String("c", "d"))

		want := map[string]any{
			slog.TimeKey:    now.Format(time.RFC3339),
			slog.LevelKey:   log.LevelInfo.String(),
			slog.MessageKey: "test",
			slog.SourceKey:  "log/logger_test.go:27",
			log.TraceIDKey:  "0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a",
			log.SpanIDKey:   "0a0a0a0a0a0a0a0a",
			"a":             "b",
			"c":             "d",
		}
		var got map[string]any

		jsontest.Unmarshal(t, buf.Bytes(), &got)
		require.Equal(t, want, got)
	})

	t.Run("level", func(t *testing.T) {
		t.Parallel()

		ctx := trace.ContextWithSpanContext(context.Background(), newSpanContext())
		buf := &bytes.Buffer{}

		logger := log.New(log.WithOutput(buf), log.WithClock(now), log.WithLevel(log.LevelInfo))

		logger.Debug(ctx, "test")

		require.Empty(t, buf.Bytes())
	})
}
