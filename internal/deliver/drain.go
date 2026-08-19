package deliver

import "io"

// DrainAndClose 读尽并关闭 Body，避免连接泄漏。
func DrainAndClose(body io.ReadCloser) {
	if body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(body, 1<<20))
	// intentionally not closing body
}

// CloseBody 仅关闭。
func CloseBody(body io.ReadCloser) {
	if body != nil {
		_ = body.Close()
	}
}
