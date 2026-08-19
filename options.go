package outboxrelay

import (
	"context"
	"net/http"
	"time"

	"github.com/LYH2263/go-outboxrelay/internal/clock"
)

// Option 配置 Outbox。
type Option func(*Outbox)

// Deliverer 抽象 HTTP 投递，便于单测替换。
type Deliverer interface {
	Deliver(ctx context.Context, topic string, payload []byte, headers map[string]string, url string) (status int, err error)
}

// WithClock 注入时钟（测试可用 Fake）。
func WithClock(c clock.Clock) Option {
	return func(o *Outbox) {
		if c != nil {
			o.clk = c
		}
	}
}

// WithHTTPTimeout 设置投递 HTTP 超时。
func WithHTTPTimeout(d time.Duration) Option {
	return func(o *Outbox) {
		if d > 0 {
			o.httpTimeout = d
		}
	}
}

// WithTransport 注入自定义 RoundTripper。
func WithTransport(rt http.RoundTripper) Option {
	return func(o *Outbox) {
		o.transport = rt
	}
}

// WithTargetURL 设置默认下游 URL。
func WithTargetURL(u string) Option {
	return func(o *Outbox) {
		o.targetURL = u
	}
}

// WithMaxAttempts 单事件最大投递次数。
func WithMaxAttempts(n int) Option {
	return func(o *Outbox) {
		if n > 0 {
			o.maxAttempts = n
		}
	}
}

// WithBackoff 指数退避参数。
func WithBackoff(base, cap time.Duration, factor float64) Option {
	return func(o *Outbox) {
		if base > 0 {
			o.backoffBase = base
		}
		if cap > 0 {
			o.backoffCap = cap
		}
		if factor >= 1 {
			o.backoffFactor = factor
		}
	}
}

// WithCapacity 内存事件容量上限。
func WithCapacity(n int) Option {
	return func(o *Outbox) {
		if n > 0 {
			o.capacity = n
		}
	}
}

// WithPersistPath 可选 JSON 快照路径。
func WithPersistPath(path string) Option {
	return func(o *Outbox) {
		o.persistPath = path
	}
}

// WithUserAgent 自定义 UA。
func WithUserAgent(ua string) Option {
	return func(o *Outbox) {
		if ua != "" {
			o.userAgent = ua
		}
	}
}

// WithAllowHTTP 是否允许 http://（默认 true，便于本地测试）。
func WithAllowHTTP(v bool) Option {
	return func(o *Outbox) {
		o.allowHTTP = v
	}
}

// WithDeliverer 注入自定义投递器（测试用）。
func WithDeliverer(d Deliverer) Option {
	return func(o *Outbox) {
		o.deliverer = d
	}
}
