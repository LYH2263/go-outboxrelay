package codec

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256Hex 载荷摘要。
func SHA256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
