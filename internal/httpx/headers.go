package httpx

import "net/http"

// CopyHeaders 拷贝。
func CopyHeaders(src http.Header) http.Header {
	if src == nil {
		return nil
	}
	dst := make(http.Header, len(src))
	for k, vs := range src {
		cp := append([]string(nil), vs...)
		dst[k] = cp
	}
	return dst
}
