package deliver

import (
	"errors"
	"fmt"
)

var (
	ErrHTTP         = errors.New("deliver: http error")
	ErrNilTransport = errors.New("deliver: nil transport")
	ErrNilClient    = errors.New("deliver: nil client")
	ErrTimeout      = errors.New("deliver: timeout")
	ErrCanceled     = errors.New("deliver: canceled")
)

// Wrap 用 %w 包裹哨兵，便于 errors.Is。
func Wrap(sentinel error, msg string) error {
	if sentinel == nil {
		return fmt.Errorf("%s", msg)
	}
	return fmt.Errorf("%w: %s", sentinel, msg)
}

// WrapErr 用 %w 链上底层错误。
func WrapErr(sentinel, err error) error {
	if err == nil {
		return sentinel
	}
	if sentinel == nil {
		return err
	}
	return fmt.Errorf("%w: %v", sentinel, err)
}
