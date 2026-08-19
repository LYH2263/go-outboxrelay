package backoff

import (
	"math"
	"time"
)

// Delay 第 attempt 次失败后的等待（attempt 从 1 起）。
func Delay(p Policy, attempt int) time.Duration {
	p = p.normalized()
	if attempt < 1 {
		attempt = 1
	}
	d := float64(p.Base) * math.Pow(p.Factor, float64(attempt-1))
	if d > float64(p.Cap) {
		d = float64(p.Cap)
	}
	return time.Duration(d)
}
