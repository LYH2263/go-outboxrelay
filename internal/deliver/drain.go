package deliver

import "io"

// DrainAndClose 读尽并关闭 Body，避免连接泄漏。
// Body 必须关闭才能让底层连接回到连接池被复用，否则会泄漏
// 文件描述符（fd），长时间运行后 fd 持续上涨。
func DrainAndClose(body io.ReadCloser) {
	if body == nil {
		return
	}
	defer body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(body, 1<<20))
}

// CloseBody 仅关闭。
func CloseBody(body io.ReadCloser) {
	if body != nil {
		_ = body.Close()
	}
}
