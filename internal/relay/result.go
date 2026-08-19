package relay

import "time"

// Outcome 单次中继结果（内部）。
type Outcome struct {
	EventID    string
	OK         bool
	StatusCode int
	Attempts   int
	Err        string
	At         time.Time
}
