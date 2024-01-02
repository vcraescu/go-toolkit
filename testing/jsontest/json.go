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
