package outboxrelay_test

import (
	"context"
	"errors"
	"testing"

	"github.com/LYH2263/go-outboxrelay"
)

func TestBug03_RelayAfterCloseNoPanic(t *testing.T) {
	o := outboxrelay.New(outboxrelay.WithTargetURL("http://127.0.0.1:9/hook"))
	if err := o.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("RelayOnce panicked after Close: %v", r)
		}
	}()
	_, err := o.RelayOnce(context.Background())
	if !errors.Is(err, outboxrelay.ErrClosed) {
		t.Fatalf("want ErrClosed, got %v", err)
	}
}
