package backoff_test

import (
	"context"
	"testing"
	"time"

	"github.com/LYH2263/go-outboxrelay/internal/backoff"
)

func TestBug08_BackoffWaitHonorsContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	err := backoff.Wait(ctx, 3*time.Second)
	if err == nil {
		t.Fatal("expected ctx error")
	}
	if time.Since(start) > time.Second {
		t.Fatalf("Wait used Sleep ignoring ctx, took %s", time.Since(start))
	}
}
