package relay

import (
	"context"
	"time"

	"github.com/LYH2263/go-outboxrelay/internal/store"
)

// BatchClaim 批量认领，最多 n 条；每轮检查 ctx。
func BatchClaim(ctx context.Context, st store.Store, now time.Time, n int) ([]store.Record, error) {
	if n <= 0 {
		n = 1
	}
	var out []store.Record
	for i := 0; i < n; i++ {
		if ctx != nil {
			if err := ctx.Err(); err != nil {
				return out, err
			}
		}
		rec, ok, err := ClaimNext(st, now)
		if err != nil {
			return out, err
		}
		if !ok {
			break
		}
		out = append(out, rec)
	}
	return out, nil
}
