package deliver

import "net/http"

// ApplyDefaults 补齐投递默认头。
func ApplyDefaults(h http.Header, topic, ua string) {
	if h == nil {
		return
	}
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", "application/json")
	}
	if ua != "" && h.Get("User-Agent") == "" {
		h.Set("User-Agent", ua)
	}
	if topic != "" {
		h.Set("X-Outbox-Topic", topic)
	}
}
