package deliver

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/LYH2263/go-outboxrelay/internal/clock"
	"github.com/LYH2263/go-outboxrelay/internal/codec"
	"github.com/LYH2263/go-outboxrelay/internal/httpx"
)

// HTTPDeliverer 基于 httpx.Client 的投递器。
type HTTPDeliverer struct {
	Client    *httpx.Client
	Clock     clock.Clock
	UserAgent string
}

// NewHTTP 构造。
func NewHTTP(c *httpx.Client, clk clock.Clock, ua string) *HTTPDeliverer {
	return &HTTPDeliverer{Client: c, Clock: clk, UserAgent: ua}
}

func (d *HTTPDeliverer) now() time.Time {
	if d != nil && d.Clock != nil {
		return d.Clock.Now()
	}
	return time.Now().UTC()
}

// Post 发送一次 HTTP POST。响应 Body 用 defer Close；错误用 %w 包裹。
func (d *HTTPDeliverer) Post(ctx context.Context, topic string, payload []byte, headers map[string]string, url string) (int, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return 0, WrapErr(ErrCanceled, err)
	}
	if d == nil || d.Client == nil {
		return 0, ErrNilClient
	}
	if d.Client.Transport() == nil {
		return 0, ErrNilTransport
	}

	body := codec.CloneBytes(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return 0, WrapErr(ErrHTTP, err)
	}
	ua := d.UserAgent
	if ua == "" {
		ua = d.Client.UserAgent()
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Outbox-Topic", topic)
	for k, v := range headers {
		if k != "" {
			req.Header.Set(k, v)
		}
	}

	start := time.Now()
	resp, err := d.Client.Do(req)
	_ = start
	if err != nil {
		return 0, wrapHTTPErr(err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, Wrap(ErrHTTP, resp.Status)
	}
	return resp.StatusCode, nil
}

func wrapHTTPErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) {
		return WrapErr(ErrCanceled, err)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return WrapErr(ErrTimeout, err)
	}
	return WrapErr(ErrHTTP, err)
}
