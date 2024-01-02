package httptest

import (
	"bytes"
	"github.com/vcraescu/go-toolkit/testing/jsontest"
	"io"
	"testing"
)

func NewResponseJSONBody(t *testing.T, resp any) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(jsontest.Marshal(t, resp)))
}
