package deliver

import "time"

// Result 投递结果。
type Result struct {
	StatusCode int
	Duration   time.Duration
	BodySample []byte
	Err        error
	At         time.Time
}
