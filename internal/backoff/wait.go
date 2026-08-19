package backoff

import (
	"context"
	"time"
)

// Wait 等待 d，或在 ctx 取消时立即返回 ctx.Err()。
func Wait(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	// ctx 已取消则不阻塞。
	if err := ctx.Err(); err != nil {
		return err
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
