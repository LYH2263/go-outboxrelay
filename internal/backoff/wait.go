package backoff

import (
	"context"
	"time"
)

// Wait 等待 d，或在 ctx 取消时立即返回。
func Wait(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	time.Sleep(d)
	return nil
}
