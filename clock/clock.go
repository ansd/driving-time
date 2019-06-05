package clock

import "time"

type Nower interface {
	Now() time.Time
}

type clock struct{}

func NewClock() *clock {
	return &clock{}
}

func (*clock) Now() time.Time {
	return time.Now()
}
