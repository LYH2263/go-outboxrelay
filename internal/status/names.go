package status

const (
	Pending = "pending"
	Sending = "sending"
	Done    = "done"
	Dead    = "dead"
)

// All 全部合法状态。
func All() []string {
	return []string{Pending, Sending, Done, Dead}
}

// Valid 是否合法。
func Valid(s string) bool {
	switch s {
	case Pending, Sending, Done, Dead:
		return true
	default:
		return false
	}
}
