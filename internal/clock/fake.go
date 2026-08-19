package clock

import (
	"sync"
	"time"
)

// Fake 可推进时钟。
type Fake struct {
	mu   sync.Mutex
	now  time.Time
}

// NewFake 从 t 开始。
func NewFake(t time.Time) *Fake {
	return &Fake{now: t.UTC()}
}

func (f *Fake) Now() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.now
}

// Advance 前进。
func (f *Fake) Advance(d time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.now = f.now.Add(d)
}

// Set 设定。
func (f *Fake) Set(t time.Time) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.now = t.UTC()
}
