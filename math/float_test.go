package math_test

import (
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/go-toolkit/math"
	"testing"
)

func TestToFixed(t *testing.T) {
	t.Parallel()

	type args struct {
		num       float64
		precision int
	}

	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "2 decimals",
			args: args{
				num:       3.5 / 2.3,
				precision: 2,
			},
			want: 1.52,
		},
		{
			name: "1 decimal",
			args: args{
				num:       1 - 0.9,
				precision: 2,
			},
			want: 0.1,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := math.ToFixed(tt.args.num, tt.args.precision)
			require.Equal(t, tt.want, got)
		})
	}
}
