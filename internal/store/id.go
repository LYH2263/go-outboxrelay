package store

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

var seq uint64

// NewID 生成可读唯一 ID。
func NewID() string {
	n := atomic.AddUint64(&seq, 1)
	var b [4]byte
	_, _ = rand.Read(b[:])
	return fmt.Sprintf("evt-%d-%s-%d", time.Now().UnixNano(), hex.EncodeToString(b[:]), n)
}
