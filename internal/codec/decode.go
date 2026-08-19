package codec

import "bytes"

// SplitParts 按 0x1e 拆分。
func SplitParts(b []byte) [][]byte {
	if len(b) == 0 {
		return nil
	}
	raw := bytes.Split(b, []byte{0x1e})
	out := make([][]byte, len(raw))
	for i, p := range raw {
		out[i] = CloneBytes(p)
	}
	return out
}
