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

func TestAppendListRelayHappy(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		defer r.Body.Close()
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	o := outboxrelay.New(
		outboxrelay.WithTargetURL(srv.URL),
		outboxrelay.WithMaxAttempts(3),
		outboxrelay.WithHTTPTimeout(2*time.Second),
	)
	defer o.Close()

	payload := []byte(`{"n":1}`)
	id, err := o.Append("order.created", payload, map[string]string{"X-A": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("empty id")
	}
	payload[0] = 'X'
	list, err := o.ListPending(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("pending=%d", len(list))
	}
	if string(list[0].Payload) != `{"n":1}` {
		t.Fatalf("payload mutated: %s", list[0].Payload)
	}
	list[0].Payload[0] = 'Y'
	again, _ := o.Get(id)
	if again.Payload[0] == 'Y' {
		t.Fatal("list alias leaked into store")
	}

	res, err := o.RelayOnce(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !res.OK || hits.Load() != 1 {
		t.Fatalf("relay=%+v hits=%d", res, hits.Load())
	}
	pend, _ := o.ListPending(10)
	if len(pend) != 0 {
		t.Fatalf("still pending: %d", len(pend))
	}
}

func TestRelayAfterClose(t *testing.T) {
	o := outboxrelay.New(outboxrelay.WithTargetURL("http://127.0.0.1:1/x"))
	_ = o.Close()
	_, err := o.RelayOnce(context.Background())
	if !errors.Is(err, outboxrelay.ErrClosed) {
		t.Fatalf("got %v", err)
	}
	_, err = o.Append("t", []byte("p"), nil)
	if !errors.Is(err, outboxrelay.ErrClosed) {
		t.Fatalf("append got %v", err)
	}
}

func TestRelayContextCancel(t *testing.T) {
	block := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-block
		w.WriteHeader(200)
	}))
	defer srv.Close()
	defer close(block)

	o := outboxrelay.New(
		outboxrelay.WithTargetURL(srv.URL),
		outboxrelay.WithHTTPTimeout(5*time.Second),
	)
	defer o.Close()
	_, _ = o.Append("t", []byte("{}"), nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := o.RelayContext(ctx, 10)
	if err == nil {
		t.Fatal("expected cancel error")
	}
}
