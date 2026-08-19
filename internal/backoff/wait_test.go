package backoff_test

import (
	"context"
	"testing"
	"time"

	"github.com/LYH2263/go-outboxrelay/internal/backoff"
)

func TestWaitHonorsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	err := backoff.Wait(ctx, time.Hour)
	if err == nil {
		t.Fatal("expected err")
	}
	if time.Since(start) > time.Second {
		t.Fatal("did not return promptly")
	}
}

func TestDelayGrows(t *testing.T) {
	p := backoff.Policy{Base: time.Millisecond, Cap: time.Second, Factor: 2}
	d1 := backoff.Delay(p, 1)
	d3 := backoff.Delay(p, 3)
	if d3 <= d1 {
		t.Fatalf("%v vs %v", d1, d3)
	}
}
