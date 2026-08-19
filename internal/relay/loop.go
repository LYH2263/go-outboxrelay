package relay

import (
	"context"
	"time"
)

// TickFunc 单次 tick。
type TickFunc func(ctx context.Context) (didWork bool, err error)

// RunLoop 在 ctx 下周期执行 tick；interval 为无工作时的休眠。
func RunLoop(ctx context.Context, interval time.Duration, tick TickFunc) error {
	if interval <= 0 {
		interval = 50 * time.Millisecond
	}
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		did, err := tick(ctx)
		if err != nil {
			return err
		}
		if did {
			continue
		}
		t := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			t.Stop()
			return ctx.Err()
		case <-t.C:
		}
	}
}
