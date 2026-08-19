package store

import (
	"sync"

	"github.com/LYH2263/go-outboxrelay/internal/clock"
)

// FileStore 基于 Memory + 强制路径的薄封装。
type FileStore struct {
	inner *Memory
	mu    sync.Mutex
}

// OpenFile 打开或创建文件快照存储。
func OpenFile(path string, capacity int, clk clock.Clock) (*FileStore, error) {
	if path == "" {
		return nil, ErrBadRecord
	}
	m := NewMemory(capacity, clk, path)
	return &FileStore{inner: m}, nil
}

func (f *FileStore) Insert(rec Record) (string, error) { return f.inner.Insert(rec) }
func (f *FileStore) Get(id string) (Record, error)       { return f.inner.Get(id) }
func (f *FileStore) ListByStatus(status string, limit int) ([]Record, error) {
	return f.inner.ListByStatus(status, limit)
}
func (f *FileStore) Update(rec Record) error               { return f.inner.Update(rec) }
func (f *FileStore) CountByStatus(status string) (int, error) {
	return f.inner.CountByStatus(status)
}
func (f *FileStore) Flush() error { return f.inner.Flush() }
func (f *FileStore) Close() error { return f.inner.Close() }
func (f *FileStore) All() ([]Record, error) { return f.inner.All() }
