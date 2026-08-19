package deliver_test

import (
	"context"
	"errors"
	"testing"

	"github.com/LYH2263/go-outboxrelay/internal/clock"
	"github.com/LYH2263/go-outboxrelay/internal/deliver"
	"github.com/LYH2263/go-outboxrelay/internal/httpx"
)

func TestBug04_NilDelivererNoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("nil deliverer panicked: %v", r)
		}
	}()
	var d *deliver.HTTPDeliverer
	_, err := d.Post(context.Background(), "t", []byte(`{}`), nil, "http://127.0.0.1:9/")
	if !errors.Is(err, deliver.ErrNilClient) {
		t.Fatalf("want ErrNilClient, got %v", err)
	}
	c := httpx.NewExact(0, nil, "ua")
	d2 := deliver.NewHTTP(c, clock.Real{}, "ua")
	_, err = d2.Post(context.Background(), "t", []byte(`{}`), nil, "http://127.0.0.1:9/")
	if !errors.Is(err, deliver.ErrNilTransport) {
		t.Fatalf("want ErrNilTransport, got %v", err)
	}
}
