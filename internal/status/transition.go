package status

import "fmt"

// CanTransition 状态机边。
func CanTransition(from, to string) bool {
	switch from {
	case Pending:
		return to == Sending || to == Dead
	case Sending:
		return to == Done || to == Pending || to == Dead
	case Done, Dead:
		return false
	default:
		return false
	}
}

// MustTransition 非法则 error。
func MustTransition(from, to string) error {
	if !CanTransition(from, to) {
		return fmt.Errorf("status: cannot %s -> %s", from, to)
	}
	return nil
}
