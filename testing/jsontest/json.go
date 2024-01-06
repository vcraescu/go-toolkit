package jsontest

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func Marshal(t *testing.T, v any) []byte {
	t.Helper()

	b, err := json.Marshal(v)
	require.NoError(t, err)

	return b
}

func Unmarshal(t *testing.T, b []byte, v any) {
	t.Helper()

	err := json.Unmarshal(b, v)
	require.NoError(t, err)
}
