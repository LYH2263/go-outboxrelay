package clock

import "time"

// Real 系统时钟。
type Real struct{}

func (Real) Now() time.Time { return time.Now().UTC() }
