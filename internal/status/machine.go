package status

import (
	"time"

	"github.com/LYH2263/go-outboxrelay/internal/backoff"
	"github.com/LYH2263/go-outboxrelay/internal/store"
)

// MarkDone 持久化成功后把事件标为 done。Update 在持久化失败时回滚内存视图，
// 故不会留下 done 的脏内存态——调用方据返回的 error 决定是否重投。
func MarkDone(st store.Store, id string) error {
	rec, err := st.Get(id)
	if err != nil {
		return err
	}
	if err := MustTransition(rec.Status, Done); err != nil {
		return err
	}
	rec.Status = Done
	rec.LastError = ""
	return st.Update(rec)
}

// MarkFail 投递失败：未超限则回 pending 并设 NextAttempt；否则 dead。
func MarkFail(st store.Store, id, msg string, pol backoff.Policy, now time.Time) (dead bool, err error) {
	rec, err := st.Get(id)
	if err != nil {
		return false, err
	}
	rec.Attempts++
	rec.LastError = msg
	if rec.Attempts >= rec.MaxAttempt {
		if err := MustTransition(rec.Status, Dead); err != nil {
			return false, err
		}
		rec.Status = Dead
		if err := st.Update(rec); err != nil {
			return false, err
		}
		return true, nil
	}
	if err := MustTransition(rec.Status, Pending); err != nil {
		return false, err
	}
	rec.Status = Pending
	rec.NextAttempt = now.Add(backoff.Delay(pol, rec.Attempts))
	if err := st.Update(rec); err != nil {
		return false, err
	}
	return false, nil
}

// MarkSending 显式进入 sending。
func MarkSending(st store.Store, id string) error {
	rec, err := st.Get(id)
	if err != nil {
		return err
	}
	if err := MustTransition(rec.Status, Sending); err != nil {
		return err
	}
	rec.Status = Sending
	return st.Update(rec)
}
