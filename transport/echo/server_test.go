package echo_test

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	kithttptest "github.com/vcraescu/go-toolkit/testing/httptest"
	kitrequire "github.com/vcraescu/go-toolkit/testing/require"
	kitecho "github.com/vcraescu/go-toolkit/transport/echo"
	"github.com/vcraescu/go-toolkit/transport/echo/handler"
	"io"
	"net/http"
	"testing"
)

func TestServer_Start(t *testing.T) {
	t.Parallel()

	type args struct {
		url    string
		body   io.Reader
		method string
	}

	tests := []struct {
		name    string
		args    args
		on      func(t *testing.T, srv *kitecho.Server)
		want    *http.Response
		wantErr bool
	}{
		{
			name: "health handler",
			args: args{
				url:    "/health",
				method: http.MethodGet,
			},
			want: &http.Response{
				StatusCode: http.StatusOK,
				Body:       kithttptest.NewResponseJSONBody(t, handler.HealthResponse{Status: "OK"}),
			},
		},
		{
			name: "new endpoint",
			args: args{
				url:    "/foo/bar",
				method: http.MethodPost,
			},
			on: func(t *testing.T, srv *kitecho.Server) {
				srv.POST("/foo/bar", func(c echo.Context) error {
					return c.JSON(http.StatusOK, map[string]string{
						"foo": "bar",
					})
				})
			},
			want: &http.Response{
				StatusCode: http.StatusOK,
				Body: kithttptest.NewResponseJSONBody(t, map[string]string{
					"foo": "bar",
				}),
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := kitecho.NewServer()

			if tt.on != nil {
				tt.on(t, srv)
			}

			client := kithttptest.StartServer(t, srv)

			req, err := http.NewRequestWithContext(context.Background(), tt.args.method, tt.args.url, tt.args.body)
			require.NoError(t, err)

			got, err := client.Do(req)
			require.NoError(t, err)

			tt.want.Header = got.Header
			kitrequire.EqualHTTPResponse(t, tt.want, got)
		})
	}
}
