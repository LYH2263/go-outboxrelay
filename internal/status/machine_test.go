package status_test

import (
	"testing"
	"time"

	"github.com/LYH2263/go-outboxrelay/internal/backoff"
	"github.com/LYH2263/go-outboxrelay/internal/clock"
	"github.com/LYH2263/go-outboxrelay/internal/status"
	"github.com/LYH2263/go-outboxrelay/internal/store"
)

func TestMarkDone(t *testing.T) {
	st := store.NewMemory(16, clock.Real{}, "")
	id, err := st.Insert(store.Record{
		Topic: "t", Payload: []byte("{}"), Status: status.Sending,
		MaxAttempt: 3, TargetURL: "http://x",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := status.MarkDone(st, id); err != nil {
		t.Fatal(err)
	}
	rec, _ := st.Get(id)
	if rec.Status != status.Done {
		t.Fatalf("%s", rec.Status)
	}
}

func TestMarkFailToPending(t *testing.T) {
	st := store.NewMemory(16, clock.Real{}, "")
	id, _ := st.Insert(store.Record{
		Topic: "t", Payload: []byte("{}"), Status: status.Sending,
		MaxAttempt: 3, TargetURL: "http://x",
	})
	dead, err := status.MarkFail(st, id, "boom", backoff.Default(), time.Now())
	if err != nil || dead {
		t.Fatalf("dead=%v err=%v", dead, err)
	}
	rec, _ := st.Get(id)
	if rec.Status != status.Pending || rec.Attempts != 1 {
		t.Fatalf("%+v", rec)
	}
}
