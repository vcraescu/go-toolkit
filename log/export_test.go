package log

import "time"

func WithClock(now time.Time) Option {
	return optionFunc(func(opts *options) {
		opts.clock = now
	})
}
