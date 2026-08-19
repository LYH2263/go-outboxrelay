package status

import "fmt"

// CheckRecord 轻量校验。
func CheckRecord(status string, attempts, max int) error {
	if !Valid(status) {
		return fmt.Errorf("status: invalid %q", status)
	}
	if attempts < 0 || max < 1 {
		return fmt.Errorf("status: bad attempts")
	}
	return nil
}
