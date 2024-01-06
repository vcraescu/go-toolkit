package require

import (
	"github.com/vcraescu/go-toolkit/testing/assert"
	"net/http"
	"testing"
)

func EqualHTTPResponse(t *testing.T, want *http.Response, got *http.Response) {
	t.Helper()

	if assert.EqualHTTPResponse(t, want, got) {
		return
	}

	t.FailNow()
}
