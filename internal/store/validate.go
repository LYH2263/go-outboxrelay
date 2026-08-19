package store

import (
	"errors"
	"strings"
)

var (
	ErrNotFound  = errors.New("store: not found")
	ErrCapacity  = errors.New("store: capacity exceeded")
	ErrClosed    = errors.New("store: closed")
	ErrBadRecord = errors.New("store: bad record")
	ErrPersist   = errors.New("store: persist failed")
)

func validateInsert(rec Record) error {
	if strings.TrimSpace(rec.Topic) == "" {
		return ErrBadRecord
	}
	if len(rec.Payload) == 0 {
		return ErrBadRecord
	}
	if strings.TrimSpace(rec.TargetURL) == "" {
		return ErrBadRecord
	}
	return nil
}
