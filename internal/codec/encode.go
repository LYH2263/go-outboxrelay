package codec

import (
	"bytes"
	"encoding/base64"
)

// EncodeB64 base64 标准编码。
func EncodeB64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// DecodeB64 解码。
func DecodeB64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// JoinParts 用 0x1e 拼接（内部信封）。
func JoinParts(parts ...[]byte) []byte {
	var buf bytes.Buffer
	for i, p := range parts {
		if i > 0 {
			buf.WriteByte(0x1e)
		}
		buf.Write(p)
	}
	return buf.Bytes()
}
