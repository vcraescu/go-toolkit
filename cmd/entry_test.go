package cmd_test

import (
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/go-toolkit/cmd"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var want bool

		handler := cmd.HandlerFunc(func(ctx cmd.Context) error {
			want = true

			return nil
		})

		cmd.Start(handler)
		require.True(t, want)
	})

	t.Run("terminate", func(t *testing.T) {
		t.Parallel()

		terminateTimeout := time.Second
		var want bool

		handler := cmd.HandlerFunc(func(ctx cmd.Context) error {
			want = true

			select {
			case <-ctx.Context().Done():
				return nil
			case <-time.NewTimer(time.Second * 2).C:
				return nil
			}
		})

		go func() {
			time.Sleep(time.Second / 3)

			cmd.Terminate()
		}()

		startTime := time.Now()
		cmd.Start(handler, cmd.WithTerminateTimeout(terminateTimeout))

		require.True(t, want)
		require.Less(t, time.Now().Sub(startTime), time.Second+time.Second/2)
	})
}
