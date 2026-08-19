package outboxrelay

import (
	"net/http"
	"sync"
	"time"

	"github.com/LYH2263/go-outboxrelay/internal/backoff"
	"github.com/LYH2263/go-outboxrelay/internal/clock"
	"github.com/LYH2263/go-outboxrelay/internal/deliver"
	"github.com/LYH2263/go-outboxrelay/internal/httpx"
	"github.com/LYH2263/go-outboxrelay/internal/store"
)

const (
	defaultHTTPTimeout = 10 * time.Second
	defaultMaxAttempts = 5
	defaultBackoffBase = 20 * time.Millisecond
	defaultBackoffCap  = 2 * time.Second
	defaultFactor      = 2.0
	defaultCapacity    = 8192
	DefaultUA          = "go-outboxrelay/1.0"
)

// Outbox 事务发件箱门面。零值不可用，须经 New 构造。
type Outbox struct {
	mu sync.Mutex

	closed bool
	clk    clock.Clock
	st     store.Store
	client *httpx.Client
	deliv  *deliver.HTTPDeliverer
	policy backoff.Policy

	httpTimeout   time.Duration
	transport     http.RoundTripper
	targetURL     string
	maxAttempts   int
	backoffBase   time.Duration
	backoffCap    time.Duration
	backoffFactor float64
	capacity      int
	persistPath   string
	userAgent     string
	allowHTTP     bool
	deliverer     Deliverer

	appended  uint64
	relayed   uint64
	succeeded uint64
	failed    uint64
	dead      uint64
	lastTopic string
	lastRelay time.Time
}

// New 构造 Outbox。
func New(opts ...Option) *Outbox {
	o := &Outbox{
		clk:           clock.Real{},
		httpTimeout:   defaultHTTPTimeout,
		maxAttempts:   defaultMaxAttempts,
		backoffBase:   defaultBackoffBase,
		backoffCap:    defaultBackoffCap,
		backoffFactor: defaultFactor,
		capacity:      defaultCapacity,
		userAgent:     DefaultUA,
		allowHTTP:     true,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
	if o.clk == nil {
		o.clk = clock.Real{}
	}
	if o.maxAttempts < 1 {
		o.maxAttempts = 1
	}
	if o.capacity < 8 {
		o.capacity = 8
	}
	o.st = store.NewMemory(o.capacity, o.clk, o.persistPath)
	o.client = httpx.New(o.httpTimeout, o.transport, o.userAgent)
	o.deliv = deliver.NewHTTP(o.client, o.clk, o.userAgent)
	o.policy = backoff.Policy{
		Base:   o.backoffBase,
		Cap:    o.backoffCap,
		Factor: o.backoffFactor,
	}
	return o
}

func (o *Outbox) checkOpenLocked() error {
	if o == nil {
		return ErrNilOutbox
	}
	if o.closed {
		return ErrClosed
	}
	return nil
}

// Store 返回底层存储（仅供内部/测试；Close 后可能为 nil）。
func (o *Outbox) Store() store.Store {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.st
}
