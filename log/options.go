package log

import (
	"io"
	"log/slog"
	"time"
)

type Level slog.Level

const (
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
	LevelDebug = Level(slog.LevelDebug)
)

func (l Level) String() string {
	return slog.Level(l).String()
}

type options struct {
	output io.Writer
	level  Level
	clock  time.Time
}

type Option interface {
	apply(opts *options)
}

var _ Option = optionFunc(nil)

type optionFunc func(opts *options)

func (fn optionFunc) apply(opts *options) {
	fn(opts)
}

func WithOutput(output io.Writer) Option {
	return optionFunc(func(opts *options) {
		opts.output = output
	})
}

func WithLevel(level Level) Option {
	return optionFunc(func(opts *options) {
		opts.level = level
	})
}
