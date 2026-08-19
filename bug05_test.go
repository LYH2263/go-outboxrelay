package outboxrelay

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-outboxrelay/internal/deliver"
)

func TestBug05_DeliverErrorWrapsSentinel(t *testing.T) {
	o := New(WithTargetURL("http://127.0.0.1:1/dead"), WithMaxAttempts(1))
	defer o.Close()
	if _, err := o.Append("order.created", []byte(`{"id":1}`), nil); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := o.RelayOnce(ctx)
	if err == nil {
		t.Fatal("expected RelayOnce delivery error")
	}
	if !errors.Is(err, deliver.ErrHTTP) && !errors.Is(err, deliver.ErrTimeout) && !errors.Is(err, deliver.ErrCanceled) {
		t.Fatalf("RelayOnce err must errors.Is deliver sentinel, got %v", err)
	}
}
