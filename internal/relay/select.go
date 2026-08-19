package relay

import (
	"sort"
	"time"

	"github.com/LYH2263/go-outboxrelay/internal/store"
)

// RankByAge 按创建时间升序。
func RankByAge(recs []store.Record) []store.Record {
	out := append([]store.Record(nil), recs...)
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out
}

// ReadyFilter 过滤到期。
func ReadyFilter(recs []store.Record, now time.Time) []store.Record {
	var out []store.Record
	for _, r := range recs {
		if !r.NextAttempt.After(now) {
			out = append(out, r)
		}
	}
	return out
}
