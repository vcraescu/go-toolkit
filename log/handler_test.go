package log_test

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/go-toolkit/log"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"strings"
	"testing"
	"time"
)

var now = time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC)

func TestNewTracedHandler(t *testing.T) {
	t.Parallel()

	type fields struct {
		opts slog.HandlerOptions
	}

	type args struct {
		ctx context.Context
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "with valid trace",
			fields: fields{},
			args: args{
				ctx: trace.ContextWithSpanContext(ctx, newSpanContext()),
			},
			want: `{"time":"2023-01-01T10:30:00Z","level":"INFO","msg":"m","a":1,"m":{"b":2},"traceID":"0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a","spanID":"0a0a0a0a0a0a0a0a"}`,
		},
		{
			name:   "with invalid valid trace",
			fields: fields{},
			args: args{
				ctx: trace.ContextWithSpanContext(ctx, trace.SpanContext{}),
			},
			want: `{"time":"2023-01-01T10:30:00Z","level":"INFO","msg":"m","a":1,"m":{"b":2}}`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			h := log.NewHandler(slog.NewJSONHandler(buf, &tt.fields.opts))

			r := slog.NewRecord(now, slog.LevelInfo, "m", 0)
			r.AddAttrs(slog.Int("a", 1), slog.Any("m", map[string]int{"b": 2}))

			err := h.Handle(tt.args.ctx, r)

			if tt.wantErr != nil {
				tt.wantErr(t, err)

				return
			}

			require.NoError(t, err)

			got := strings.TrimSuffix(buf.String(), "\n")
			require.Equal(t, tt.want, got)
		})
	}
}

func newSpanContext() trace.SpanContext {
	traceID := trace.TraceID{}

	for i := range traceID {
		traceID[i] = 10
	}

	spanID := trace.SpanID{}

	for i := range spanID {
		spanID[i] = 10
	}

	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
	})
}
