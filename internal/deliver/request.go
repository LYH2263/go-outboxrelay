package deliver

import "github.com/LYH2263/go-outboxrelay/internal/codec"

// Request HTTP 投递请求。
type Request struct {
	Topic   string
	URL     string
	Payload []byte
	Headers map[string]string
}

func cloneRequest(in Request) Request {
	return Request{
		Topic:   in.Topic,
		URL:     in.URL,
		Payload: codec.CloneBytes(in.Payload),
		Headers: codec.CloneHeaders(in.Headers),
	}
}
