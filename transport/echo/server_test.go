package echo_test

import (
	"context"
	"github.com/stretchr/testify/require"
	kithttptest "github.com/vcraescu/go-toolkit/testing/httptest"
	kitrequire "github.com/vcraescu/go-toolkit/testing/require"
	"github.com/vcraescu/go-toolkit/transport/echo"
	"github.com/vcraescu/go-toolkit/transport/echo/handler"
	"io"
	"net/http"
	"testing"
)

func TestServer_Start(t *testing.T) {
	t.Parallel()

	type args struct {
		url  string
		body io.Reader
	}

	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "health handler",
			args: args{
				url: "/health",
			},
			want: &http.Response{
				StatusCode: http.StatusOK,
				Body:       kithttptest.NewResponseJSONBody(t, handler.HealthResponse{Status: "OK"}),
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := kithttptest.StartServer(t, echo.NewServer())

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, tt.args.url, tt.args.body)
			require.NoError(t, err)

			got, err := client.Do(req)
			require.NoError(t, err)

			tt.want.Header = got.Header
			kitrequire.EqualHTTPResponse(t, tt.want, got)
		})
	}
}
