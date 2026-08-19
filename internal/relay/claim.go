package relay

import (
	"time"

	"github.com/LYH2263/go-outboxrelay/internal/store"
)

// ClaimNext 取下一条到期 pending 并标为 sending。
func ClaimNext(st store.Store, now time.Time) (store.Record, bool, error) {
	if st == nil {
		return store.Record{}, false, store.ErrClosed
	}
	var due []store.Record
	if m, ok := st.(*store.Memory); ok {
		var err error
		due, err = m.DuePending(now, 1)
		if err != nil {
			return store.Record{}, false, err
		}
	} else {
		recs, err := st.ListByStatus("pending", 64)
		if err != nil {
			return store.Record{}, false, err
		}
		for _, r := range recs {
			if !r.NextAttempt.After(now) {
				due = append(due, r)
				break
			}
		}
	}
	if len(due) == 0 {
		return store.Record{}, false, nil
	}
	rec, err := store.ClaimSending(st, due[0].ID)
	if err != nil {
		return store.Record{}, false, err
	}
	return rec, true, nil
}
