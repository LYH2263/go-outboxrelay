package outboxrelay_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/LYH2263/go-outboxrelay"
)

func TestBug07_RelayContextHonorsCancel(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		defer r.Body.Close()
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	o := outboxrelay.New(
		outboxrelay.WithTargetURL(srv.URL),
		outboxrelay.WithHTTPTimeout(2*time.Second),
	)
	defer o.Close()
	for i := 0; i < 5; i++ {
		if _, err := o.Append("t", []byte(`{}`), nil); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	res, err := o.RelayContext(ctx, 5)
	if err == nil {
		t.Fatal("expected cancel error")
	}
	if !(errors.Is(err, outboxrelay.ErrCanceled) || errors.Is(err, context.Canceled)) {
		t.Fatalf("want canceled, got %v", err)
	}
	if hits.Load() != 0 {
		t.Fatalf("RelayContext ignored cancel and delivered hits=%d results=%d", hits.Load(), len(res))
	}
	if time.Since(start) > time.Second {
		t.Fatalf("RelayContext ignored cancel, took %s", time.Since(start))
	}
}
