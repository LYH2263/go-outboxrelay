package codec

// CloneBytes 返回 payload 的独立副本。
// 必须分配新底层数组，否则调用方复用 reusable []byte 缓冲区
// （原地改写或 append 追加 trace id）会污染已入库记录。
func CloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	out := make([]byte, len(b))
	copy(out, b)
	return out
}

// CloneHeaders 深拷贝 header map。
func CloneHeaders(h map[string]string) map[string]string {
	if h == nil {
		return nil
	}
	out := make(map[string]string, len(h))
	for k, v := range h {
		out[k] = v
	}
	return out
}
