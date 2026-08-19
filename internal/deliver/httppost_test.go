package deliver_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LYH2263/go-outboxrelay/internal/clock"
	"github.com/LYH2263/go-outboxrelay/internal/deliver"
	"github.com/LYH2263/go-outboxrelay/internal/httpx"
)

func TestPostOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(200)
	}))
	defer srv.Close()
	c := httpx.New(0, nil, "test")
	d := deliver.NewHTTP(c, clock.Real{}, "test")
	code, err := d.Post(context.Background(), "t", []byte("{}"), nil, srv.URL)
	if err != nil || code != 200 {
		t.Fatalf("%d %v", code, err)
	}
}

func TestPostWraps(t *testing.T) {
	c := httpx.New(0, nil, "test")
	d := deliver.NewHTTP(c, clock.Real{}, "test")
	_, err := d.Post(context.Background(), "t", []byte("{}"), nil, "http://127.0.0.1:1/")
	if err == nil || !errors.Is(err, deliver.ErrHTTP) {
		t.Fatalf("got %v", err)
	}
}

// TestPostCanceledBefore 不发请求，立即返回可识别的取消错误。
func TestPostCanceledBefore(t *testing.T) {
	c := httpx.New(0, nil, "test")
	d := deliver.NewHTTP(c, clock.Real{}, "test")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	_, err := d.Post(ctx, "t", []byte("{}"), nil, "http://127.0.0.1:1/")
	if err == nil {
		t.Fatal("expected cancel error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, got %v", err)
	}
	if !errors.Is(err, deliver.ErrCanceled) {
		t.Fatalf("want deliver.ErrCanceled, got %v", err)
	}
	if time.Since(start) > time.Second {
		t.Fatalf("Post ignored pre-canceled ctx, took %s", time.Since(start))
	}
}

// TestPostCanceledMidFlight 投递过程中取消：Do 应因 ctx 立即中止，
// 而非等到 http.Client.Timeout 或下游解除阻塞。
func TestPostCanceledMidFlight(t *testing.T) {
	block := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-block
		w.WriteHeader(200)
	}))
	defer srv.Close()
	defer close(block)

	c := httpx.New(0, nil, "test") // 默认 10s 客户端超时
	d := deliver.NewHTTP(c, clock.Real{}, "test")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	start := time.Now()
	_, err := d.Post(ctx, "t", []byte("{}"), nil, srv.URL)
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("expected cancel error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, got %v", err)
	}
	if !errors.Is(err, deliver.ErrCanceled) {
		t.Fatalf("want deliver.ErrCanceled, got %v", err)
	}
	if elapsed > time.Second {
		t.Fatalf("Post did not honor in-flight cancel, took %s", elapsed)
	}
}
