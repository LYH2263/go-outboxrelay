package httpx

// IsSuccess 2xx。
func IsSuccess(code int) bool {
	return code >= 200 && code < 300
}

// IsRetryable 5xx 或 429。
func IsRetryable(code int) bool {
	return code == 429 || code >= 500
}
