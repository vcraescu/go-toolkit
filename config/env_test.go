package config_test

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/go-toolkit/config"
	"strconv"
	"testing"
)

func Test_env_Get(t *testing.T) {
	type args struct {
		key  string
		opts []config.GetOption
	}

	tests := []struct {
		name    string
		args    args
		on      func(t *testing.T)
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty key",
			args: args{},
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, config.ErrEmptyKey)
			},
		},
		{
			name: "not found",
			args: args{
				key: "FOOBAR",
			},
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, config.ErrNotFound)
			},
		},
		{
			name: "simple string",
			args: args{
				key: "FOO",
			},
			on: func(t *testing.T) {
				t.Setenv("FOO", "BAR")
			},
			want: "BAR",
		},
		{
			name: "base64 encoded string",
			args: args{
				key: "FOO",
				opts: []config.GetOption{
					config.WithBase64(),
				},
			},
			on: func(t *testing.T) {
				t.Setenv("FOO", base64.StdEncoding.EncodeToString([]byte("BAR")))
			},
			want: "BAR",
		},
		{
			name: "gzip encoded string",
			args: args{
				key: "FOO",
				opts: []config.GetOption{
					config.WithGZIP(),
				},
			},
			on: func(t *testing.T) {
				t.Setenv("FOO", mustGZIPEncode(t, "BAR"))
			},
			want: "BAR",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if tt.on != nil {
				tt.on(t)
			}

			cfg := config.NewFromEnv()
			got, err := cfg.Get(tt.args.key, tt.args.opts...)

			if tt.wantErr != nil {
				tt.wantErr(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_env_GetFloat(t *testing.T) {
	type args struct {
		key  string
		opts []config.GetOption
	}

	tests := []struct {
		name    string
		args    args
		on      func(t *testing.T)
		want    float64
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty key",
			args: args{},
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, config.ErrEmptyKey)
			},
		},
		{
			name: "not found",
			args: args{
				key: "FOOBAR",
			},
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, config.ErrNotFound)
			},
		},
		{
			name: "non float",
			args: args{
				key: "FOO",
			},
			on: func(t *testing.T) {
				t.Setenv("FOO", "BAR")
			},
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, strconv.ErrSyntax)
			},
		},
		{
			name: "valid float",
			args: args{
				key: "FOO",
			},
			on: func(t *testing.T) {
				t.Setenv("FOO", "12.5")
			},
			want: 12.5,
		},
		{
			name: "encoded valid float",
			args: args{
				key: "FOO",
				opts: []config.GetOption{
					config.WithBase64(),
				},
			},
			on: func(t *testing.T) {
				t.Setenv("FOO", base64.StdEncoding.EncodeToString([]byte("12.5")))
			},
			want: 12.5,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if tt.on != nil {
				tt.on(t)
			}

			cfg := config.NewFromEnv()
			got, err := cfg.GetFloat(tt.args.key, tt.args.opts...)

			if tt.wantErr != nil {
				tt.wantErr(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_env_Have(t *testing.T) {
	type args struct {
		keys []string
	}

	tests := []struct {
		name    string
		args    args
		on      func(t *testing.T)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "missing keys",
			args: args{
				keys: []string{"FOO", "BAR"},
			},
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, config.ErrNotFound)
			},
		},
		{
			name: "non missing keys",
			args: args{
				keys: []string{"FOO", "BAR"},
			},
			on: func(t *testing.T) {
				t.Setenv("FOO", "")
				t.Setenv("BAR", "")
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if tt.on != nil {
				tt.on(t)
			}

			err := config.NewFromEnv().Have(tt.args.keys...)

			if tt.wantErr != nil {
				tt.wantErr(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

func mustGZIPEncode(t *testing.T, s string) string {
	t.Helper()

	buf := &bytes.Buffer{}
	w := gzip.NewWriter(buf)

	_, err := w.Write([]byte(s))
	require.NoError(t, err)

	err = w.Close()
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
