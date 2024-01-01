package timeclock

import (
	"time"
)

var (
	_   Clock = Func(nil)
	Now Clock = Func(time.Now)
)

type Clock interface {
	Now() time.Time
}

type Func func() time.Time

func (fn Func) Now() time.Time {
	return fn()
}
