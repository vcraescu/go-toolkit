package httptest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func StartServer(t *testing.T, handler http.Handler) *http.Client {
	t.Helper()

	srv := httptest.NewUnstartedServer(handler)
	srv.Start()

	t.Cleanup(func() {
		srv.Close()
	})

	return &http.Client{
		Transport: NewRoundTripperWithBaseURL(srv.URL),
	}
}

type RoundTripper struct {
	baseURL   string
	transport http.RoundTripper
	t         *testing.T
}

func NewRoundTripperWithBaseURL(baseURL string) *RoundTripper {
	return &RoundTripper{
		baseURL:   baseURL,
		transport: http.DefaultTransport,
	}
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	u, err := url.Parse(r.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse baseURL: %w", err)
	}

	req.URL.Host = u.Host
	req.URL.Scheme = u.Scheme

	return r.transport.RoundTrip(req)
}
