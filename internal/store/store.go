package store

import "time"

// Store 发件箱持久化抽象。
type Store interface {
	Insert(rec Record) (id string, err error)
	Get(id string) (Record, error)
	ListByStatus(status string, limit int) ([]Record, error)
	Update(rec Record) error
	CountByStatus(status string) (int, error)
	Flush() error
	Close() error
	All() ([]Record, error)
}

// Record 内部存储记录。Payload/Headers 由调用方保证所有权语义。
type Record struct {
	ID          string
	Topic       string
	Payload     []byte
	Headers     map[string]string
	Status      string
	Attempts    int
	MaxAttempt  int
	LastError   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	NextAttempt time.Time
	TargetURL   string
}
