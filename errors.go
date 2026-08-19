package outboxrelay

import "errors"

var (
	ErrClosed       = errors.New("outboxrelay: closed")
	ErrNilOutbox    = errors.New("outboxrelay: nil outbox")
	ErrEmptyTopic   = errors.New("outboxrelay: empty topic")
	ErrEmptyPayload = errors.New("outboxrelay: empty payload")
	ErrEmptyURL     = errors.New("outboxrelay: empty target url")
	ErrNotFound     = errors.New("outboxrelay: event not found")
	ErrInvalidStatus = errors.New("outboxrelay: invalid status")
	ErrPersist      = errors.New("outboxrelay: persist failed")
	ErrDeliver      = errors.New("outboxrelay: deliver failed")
	ErrHTTP         = errors.New("outboxrelay: http error")
	ErrTimeout      = errors.New("outboxrelay: timeout")
	ErrCanceled     = errors.New("outboxrelay: canceled")
	ErrNilTransport = errors.New("outboxrelay: nil transport")
	ErrCapacity     = errors.New("outboxrelay: capacity exceeded")
)
