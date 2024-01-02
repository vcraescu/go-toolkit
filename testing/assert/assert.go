package assert

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func EqualHTTPResponse(t *testing.T, want *http.Response, got *http.Response) bool {
	t.Helper()

	if !assert.Equal(t, want.StatusCode, got.StatusCode) {
		return false
	}

	if want.Body == got.Body {
		return true
	}

	var wantBody, gotBody []byte

	var err error

	if want.Body != nil {
		wantBody, err = io.ReadAll(want.Body)
		if !assert.NoError(t, err) {
			return false
		}

		wantBody = bytes.TrimSpace(wantBody)
	}

	if got.Body != nil {
		gotBody, err = io.ReadAll(got.Body)
		if !assert.NoError(t, err) {
			return false
		}

		gotBody = bytes.TrimSpace(gotBody)
	}

	if !assert.Equal(t, string(wantBody), string(gotBody)) {
		return false
	}

	if !assert.Equal(t, want.Header, got.Header) {
		return false
	}

	return true
}
