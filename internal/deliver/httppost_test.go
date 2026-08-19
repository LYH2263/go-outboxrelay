package deliver_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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
