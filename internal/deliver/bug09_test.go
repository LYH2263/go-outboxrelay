package deliver_test

import (
	"context"
	"io"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/LYH2263/go-outboxrelay/internal/clock"
	"github.com/LYH2263/go-outboxrelay/internal/deliver"
	"github.com/LYH2263/go-outboxrelay/internal/httpx"
)

type closeCounter struct {
	n atomic.Int32
}

func (c *closeCounter) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       &countCloser{c: c, r: http.NoBody},
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

type countCloser struct {
	c *closeCounter
	r io.ReadCloser
}

func (c *countCloser) Read(p []byte) (int, error) { return c.r.Read(p) }
func (c *countCloser) Close() error {
	c.c.n.Add(1)
	return c.r.Close()
}

func TestBug09_ResponseBodyClosed(t *testing.T) {
	ctr := &closeCounter{}
	client := httpx.NewExact(0, ctr, "ua")
	d := deliver.NewHTTP(client, clock.Real{}, "ua")
	_, err := d.Post(context.Background(), "t", []byte(`{}`), nil, "http://example.invalid/hook")
	if err != nil {
		t.Fatal(err)
	}
	if ctr.n.Load() < 1 {
		t.Fatalf("response Body Close count=%d, want >=1", ctr.n.Load())
	}
}
