package deliver

import (
	"net/url"
	"strings"
)

// ValidateURL 检查投递 URL。
func ValidateURL(raw string, allowHTTP bool) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Wrap(ErrHTTP, "empty url")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return WrapErr(ErrHTTP, err)
	}
	switch strings.ToLower(u.Scheme) {
	case "https":
		return nil
	case "http":
		if allowHTTP {
			return nil
		}
		return Wrap(ErrHTTP, "http not allowed")
	default:
		return Wrap(ErrHTTP, "unsupported scheme")
	}
}
