package store

import (
	"sort"
	"sync"
	"time"

	"github.com/LYH2263/go-outboxrelay/internal/clock"
)

// Memory 线程安全内存存储，可选 JSON 快照。
type Memory struct {
	mu       sync.Mutex
	cap      int
	clk      clock.Clock
	path     string
	closed   bool
	byID     map[string]*Record
	order    []string
	dirty    bool
	flushErr error // 测试可注入
}

// NewMemory 构造内存存储。
func NewMemory(capacity int, clk clock.Clock, persistPath string) *Memory {
	if capacity < 1 {
		capacity = 64
	}
	if clk == nil {
		clk = clock.Real{}
	}
	m := &Memory{
		cap:  capacity,
		clk:  clk,
		path: persistPath,
		byID: make(map[string]*Record),
	}
	if persistPath != "" {
		_ = m.loadLocked()
	}
	return m
}

// SetFlushError 测试注入：下次 Flush 或带路径的 persist 返回该错误。
func (m *Memory) SetFlushError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.flushErr = err
}

func (m *Memory) Insert(rec Record) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return "", ErrClosed
	}
	if err := validateInsert(rec); err != nil {
		return "", err
	}
	if len(m.byID) >= m.cap {
		return "", ErrCapacity
	}
	now := m.clk.Now()
	id := rec.ID
	if id == "" {
		id = NewID()
	}
	cp := CloneRecord(rec)
	cp.ID = id
	cp.CreatedAt = now
	cp.UpdatedAt = now
	if cp.NextAttempt.IsZero() {
		cp.NextAttempt = now
	}
	m.byID[id] = &cp
	m.order = append(m.order, id)
	m.dirty = true
	if m.path != "" {
		if err := m.persistLocked(); err != nil {
			delete(m.byID, id)
			m.order = m.order[:len(m.order)-1]
			return "", err
		}
	}
	return id, nil
}

func (m *Memory) Get(id string) (Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return Record{}, ErrClosed
	}
	r, ok := m.byID[id]
	if !ok {
		return Record{}, ErrNotFound
	}
	return CloneRecord(*r), nil
}

func (m *Memory) ListByStatus(status string, limit int) ([]Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil, ErrClosed
	}
	var out []Record
	for _, id := range m.order {
		r := m.byID[id]
		if r == nil || r.Status != status {
			continue
		}
		out = append(out, CloneRecord(*r))
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (m *Memory) Update(rec Record) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return ErrClosed
	}
	cur, ok := m.byID[rec.ID]
	if !ok {
		return ErrNotFound
	}
	cp := CloneRecord(rec)
	cp.CreatedAt = cur.CreatedAt
	cp.UpdatedAt = m.clk.Now()
	// 先尝试持久化，成功后再写内存
	prev := *cur
	*cur = cp
	m.dirty = true
	if m.path != "" {
		if err := m.persistLocked(); err != nil {
			*cur = prev
			m.dirty = true
			return err
		}
	}
	return nil
}

func (m *Memory) CountByStatus(status string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return 0, ErrClosed
	}
	n := 0
	for _, r := range m.byID {
		if r.Status == status {
			n++
		}
	}
	return n, nil
}

func (m *Memory) All() ([]Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil, ErrClosed
	}
	out := make([]Record, 0, len(m.order))
	for _, id := range m.order {
		if r := m.byID[id]; r != nil {
			out = append(out, CloneRecord(*r))
		}
	}
	return out, nil
}

func (m *Memory) Flush() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return ErrClosed
	}
	if m.flushErr != nil {
		err := m.flushErr
		m.flushErr = nil
		return err
	}
	if m.path == "" {
		m.dirty = false
		return nil
	}
	if err := m.persistLocked(); err != nil {
		return err
	}
	m.dirty = false
	return nil
}

func (m *Memory) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil
	}
	m.closed = true
	m.byID = nil
	m.order = nil
	m.path = ""
	return nil
}

// DuePending 返回到期可投递的 pending（按 NextAttempt、创建At）。
func (m *Memory) DuePending(now time.Time, limit int) ([]Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil, ErrClosed
	}
	type pair struct {
		id string
		t  time.Time
	}
	var cand []pair
	for _, id := range m.order {
		r := m.byID[id]
		if r == nil || r.Status != "pending" {
			continue
		}
		if r.NextAttempt.After(now) {
			continue
		}
		cand = append(cand, pair{id: id, t: r.CreatedAt})
	}
	sort.Slice(cand, func(i, j int) bool {
		return cand[i].t.Before(cand[j].t)
	})
	var out []Record
	for _, c := range cand {
		out = append(out, CloneRecord(*m.byID[c.id]))
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}
