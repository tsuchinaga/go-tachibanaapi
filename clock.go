package tachibana

import "time"

// newClock - iClockを生成する
func newClock() iClock {
	return &clock{}
}

// iClock - 時間関連のインターフェース
type iClock interface {
	Now() time.Time
}

// clock - 時間関連
type clock struct{}

func (c *clock) Now() time.Time {
	return time.Now()
}
