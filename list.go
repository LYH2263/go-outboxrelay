package outboxrelay

import (
	"github.com/LYH2263/go-outboxrelay/internal/codec"
	"github.com/LYH2263/go-outboxrelay/internal/status"
	"github.com/LYH2263/go-outboxrelay/internal/store"
)

// ListPending 返回待投递事件的深拷贝视图（切片与 payload/headers 均独立）。
func (o *Outbox) ListPending(limit int) ([]Event, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if err := o.checkOpenLocked(); err != nil {
		return nil, err
	}
	recs, err := o.st.ListByStatus(string(status.Pending), limit)
	if err != nil {
		return nil, err
	}
	out := make([]Event, 0, len(recs))
	for _, r := range recs {
		out = append(out, recordToEvent(r))
	}
	return out, nil
}

// ListByStatus 按状态列出事件（深拷贝）。
func (o *Outbox) ListByStatus(st Status, limit int) ([]Event, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if err := o.checkOpenLocked(); err != nil {
		return nil, err
	}
	recs, err := o.st.ListByStatus(string(st), limit)
	if err != nil {
		return nil, err
	}
	out := make([]Event, 0, len(recs))
	for _, r := range recs {
		out = append(out, recordToEvent(r))
	}
	return out, nil
}

// Get 按 ID 取事件（深拷贝）。
func (o *Outbox) Get(id string) (Event, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if err := o.checkOpenLocked(); err != nil {
		return Event{}, err
	}
	r, err := o.st.Get(id)
	if err != nil {
		return Event{}, err
	}
	return recordToEvent(r), nil
}

func recordToEvent(r store.Record) Event {
	return Event{
		ID:          r.ID,
		Topic:       r.Topic,
		Payload:     r.Payload,
		Headers:     r.Headers,
		Status:      Status(r.Status),
		Attempts:    r.Attempts,
		MaxAttempt:  r.MaxAttempt,
		LastError:   r.LastError,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		NextAttempt: r.NextAttempt,
		TargetURL:   r.TargetURL,
	}
}
