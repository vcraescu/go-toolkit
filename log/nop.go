package log

import "context"

var _ Logger = (*Nop)(nil)

type Nop struct{}

func (l Nop) With(args ...any) Logger {
	return l
}

func NewNop() *Nop {
	return &Nop{}
}

func (l Nop) Info(ctx context.Context, msg string, args ...any) {}

func (l Nop) Error(ctx context.Context, msg string, args ...any) {}

func (l Nop) Warn(ctx context.Context, msg string, args ...any) {}

func (l Nop) Debug(ctx context.Context, msg string, args ...any) {}
