package outboxrelay_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LYH2263/go-outboxrelay"
	"github.com/LYH2263/go-outboxrelay/internal/store"
)

func TestBug10_CloseFlushesStoreBeforeRelease(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "outbox.json")
	o := outboxrelay.New(
		outboxrelay.WithTargetURL("http://127.0.0.1:9/hook"),
		outboxrelay.WithPersistPath(path),
	)
	id, err := o.Append("order.created", []byte(`{"n":1}`), nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := o.Close(); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) == 0 {
		t.Fatal("persist file empty after Close")
	}
	recs, err := store.LoadJSON(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) != 1 || recs[0].ID != id {
		t.Fatalf("Close released store before Flush; loaded=%+v", recs)
	}
}
