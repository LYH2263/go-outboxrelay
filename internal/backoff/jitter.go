package backoff

import (
	"math/rand"
	"time"
)

// WithJitter 在 [0.5d, 1.5d] 加抖动。
func WithJitter(d time.Duration, rng *rand.Rand) time.Duration {
	if d <= 0 {
		return 0
	}
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	f := 0.5 + rng.Float64()
	return time.Duration(float64(d) * f)
}
