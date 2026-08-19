package codec

import (
	"encoding/json"
	"time"
)

// Envelope 投递信封。
type Envelope struct {
	Topic     string            `json:"topic"`
	Payload   json.RawMessage   `json:"payload"`
	Headers   map[string]string `json:"headers,omitempty"`
	Timestamp time.Time         `json:"ts"`
}

// WrapEnvelope 包装。
func WrapEnvelope(topic string, payload []byte, headers map[string]string, ts time.Time) ([]byte, error) {
	env := Envelope{
		Topic:     topic,
		Payload:   json.RawMessage(CloneBytes(payload)),
		Headers:   CloneHeaders(headers),
		Timestamp: ts.UTC(),
	}
	return json.Marshal(env)
}
