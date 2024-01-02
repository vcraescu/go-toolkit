package log_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/go-toolkit/log"
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

		logger := log.New(log.WithOutput(buf), log.WithClock(now))

		logger.Info(ctx, "test", slog.String("c", "d"))

		want := map[string]string{
			slog.TimeKey:    now.Format(time.RFC3339),
			slog.LevelKey:   log.LevelInfo.String(),
			slog.MessageKey: "test",
			slog.SourceKey:  "log/logger_test.go:26",
			log.TraceIDKey:  "0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a",
			log.SpanIDKey:   "0a0a0a0a0a0a0a0a",
			"c":             "d",
		}
		var got map[string]string

		mustUnmarshal(t, buf.Bytes(), &got)
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

func mustUnmarshal(t *testing.T, b []byte, v any) {
	err := json.Unmarshal(b, v)
	require.NoError(t, err)
}
