package outboxrelay

import "time"

// Status 发件箱事件生命周期。
type Status string

const (
	StatusPending Status = "pending"
	StatusSending Status = "sending"
	StatusDone    Status = "done"
	StatusDead    Status = "dead"
)

// Event 对外可见的发件箱事件视图（载荷与头已拷贝）。
type Event struct {
	ID         string
	Topic      string
	Payload    []byte
	Headers    map[string]string
	Status     Status
	Attempts   int
	MaxAttempt int
	LastError  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	NextAttempt time.Time
	TargetURL  string
}

// RelayResult 单次投递结果。
type RelayResult struct {
	EventID    string
	Topic      string
	OK         bool
	StatusCode int
	Attempts   int
	Duration   time.Duration
	Err        string
	FinishedAt time.Time
}

// Stats 运行统计快照。
type Stats struct {
	Appended   uint64
	Relayed    uint64
	Succeeded  uint64
	Failed     uint64
	Dead       uint64
	Pending    int
	Closed     bool
	LastTopic  string
	LastRelay  time.Time
}
