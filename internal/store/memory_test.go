package store_test

import (
	"testing"

	"github.com/LYH2263/go-outboxrelay/internal/clock"
	"github.com/LYH2263/go-outboxrelay/internal/store"
)

func TestMemoryInsertGet(t *testing.T) {
	m := store.NewMemory(8, clock.Real{}, "")
	id, err := m.Insert(store.Record{
		Topic: "a", Payload: []byte("p"), Status: "pending", TargetURL: "http://x", MaxAttempt: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	rec, err := m.Get(id)
	if err != nil || string(rec.Payload) != "p" {
		t.Fatalf("%+v %v", rec, err)
	}
	rec.Payload[0] = 'Z'
	again, _ := m.Get(id)
	if again.Payload[0] == 'Z' {
		t.Fatal("get alias")
	}
}

func TestFlushCloseOrder(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/snap.json"
	m := store.NewMemory(8, clock.Real{}, path)
	_, _ = m.Insert(store.Record{
		Topic: "a", Payload: []byte("p"), Status: "pending", TargetURL: "http://x", MaxAttempt: 2,
	})
	if err := m.Flush(); err != nil {
		t.Fatal(err)
	}
	if err := m.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Get("x"); err == nil {
		t.Fatal("expected closed")
	}
}
