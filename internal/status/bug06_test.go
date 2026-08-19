package status_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/LYH2263/go-outboxrelay/internal/clock"
	"github.com/LYH2263/go-outboxrelay/internal/status"
	"github.com/LYH2263/go-outboxrelay/internal/store"
)

func TestBug06_MarkDonePersistFailureRollsBack(t *testing.T) {
	path := filepath.Join(t.TempDir(), "snap.json")
	st := store.NewMemory(16, clock.Real{}, path)
	id, err := st.Insert(store.Record{
		Topic: "t", Payload: []byte("{}"), Status: status.Sending,
		MaxAttempt: 3, TargetURL: "http://x",
	})
	if err != nil {
		t.Fatal(err)
	}
	st.SetFlushError(errors.New("disk full"))
	err = status.MarkDone(st, id)
	if err == nil {
		t.Fatal("expected persist error")
	}
	rec, err := st.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Status == status.Done {
		t.Fatalf("memory mutated to done despite persist failure: %+v", rec)
	}
	if rec.Status != status.Sending {
		t.Fatalf("want sending rollback, got %s", rec.Status)
	}
}
