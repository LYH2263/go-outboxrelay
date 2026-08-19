package outboxrelay

import (
	"strings"

	"github.com/LYH2263/go-outboxrelay/internal/codec"
	"github.com/LYH2263/go-outboxrelay/internal/status"
	"github.com/LYH2263/go-outboxrelay/internal/store"
)

// Append 将事件写入发件箱。payload 与 headers 会被深拷贝，调用方之后修改不影响库存。
func (o *Outbox) Append(topic string, payload []byte, headers map[string]string) (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if err := o.checkOpenLocked(); err != nil {
		return "", err
	}
	topic = strings.TrimSpace(topic)
	if topic == "" {
		return "", ErrEmptyTopic
	}
	if len(payload) == 0 {
		return "", ErrEmptyPayload
	}
	url := o.targetURL
	if url == "" {
		return "", ErrEmptyURL
	}
	rec := store.Record{
		Topic:      topic,
		Payload:    payload,
		Headers:    codec.CloneHeaders(headers),
		Status:     string(status.Pending),
		MaxAttempt: o.maxAttempts,
		TargetURL:  url,
	}
	id, err := o.st.Insert(rec)
	if err != nil {
		return "", err
	}
	o.appended++
	o.lastTopic = topic
	return id, nil
}

// AppendEvent 使用完整 Event 字段追加（忽略 ID/Status，由库生成）。
func (o *Outbox) AppendEvent(ev Event) (string, error) {
	headers := ev.Headers
	if headers == nil {
		headers = map[string]string{}
	}
	if ev.TargetURL != "" {
		o.mu.Lock()
		prev := o.targetURL
		o.targetURL = ev.TargetURL
		o.mu.Unlock()
		id, err := o.Append(ev.Topic, ev.Payload, headers)
		o.mu.Lock()
		o.targetURL = prev
		o.mu.Unlock()
		return id, err
	}
	return o.Append(ev.Topic, ev.Payload, headers)
}
