package outboxrelay

import "github.com/LYH2263/go-outboxrelay/internal/status"

// Snapshot 返回统计快照。
func (o *Outbox) Snapshot() Stats {
	o.mu.Lock()
	defer o.mu.Unlock()
	s := Stats{
		Appended:  o.appended,
		Relayed:   o.relayed,
		Succeeded: o.succeeded,
		Failed:    o.failed,
		Dead:      o.dead,
		Closed:    o.closed,
		LastTopic: o.lastTopic,
		LastRelay: o.lastRelay,
	}
	if o.st != nil {
		if n, err := o.st.CountByStatus(string(status.Pending)); err == nil {
			s.Pending = n
		}
	}
	return s
}
