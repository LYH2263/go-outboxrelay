package httpx

import (
	"net/http"
	"time"
)

// Client 薄封装 http.Client，可查询 Transport。
type Client struct {
	hc        *http.Client
	ua        string
	transport http.RoundTripper
}

// New 构造。transport 为 nil 时使用 http.DefaultTransport。
func New(timeout time.Duration, transport http.RoundTripper, ua string) *Client {
	rt := transport
	if rt == nil {
		rt = http.DefaultTransport
	}
	return newClient(timeout, rt, ua)
}

// NewExact 不替换 nil transport（供防御性路径测试）。
func NewExact(timeout time.Duration, transport http.RoundTripper, ua string) *Client {
	return newClient(timeout, transport, ua)
}

func newClient(timeout time.Duration, rt http.RoundTripper, ua string) *Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Client{
		hc: &http.Client{
			Timeout:   timeout,
			Transport: rt,
		},
		ua:        ua,
		transport: rt,
	}
}

// Do 发请求。
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c == nil || c.hc == nil {
		return nil, errNil
	}
	return c.hc.Do(req)
}

// Transport 返回底层 RoundTripper。
func (c *Client) Transport() http.RoundTripper {
	if c == nil {
		return nil
	}
	return c.transport
}

// UserAgent UA。
func (c *Client) UserAgent() string {
	if c == nil {
		return ""
	}
	return c.ua
}

// CloseIdle 关闭空闲连接。
func (c *Client) CloseIdle() {
	if c == nil || c.hc == nil {
		return
	}
	c.hc.CloseIdleConnections()
}

type nilErr string

func (e nilErr) Error() string { return string(e) }

const errNil = nilErr("httpx: nil client")
