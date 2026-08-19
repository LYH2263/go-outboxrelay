package store

// ClaimSending 将 pending 标记为 sending（内存乐观更新）。
func ClaimSending(st Store, id string) (Record, error) {
	rec, err := st.Get(id)
	if err != nil {
		return Record{}, err
	}
	if rec.Status != "pending" {
		return Record{}, ErrBadRecord
	}
	rec.Status = "sending"
	if err := st.Update(rec); err != nil {
		return Record{}, err
	}
	return st.Get(id)
}

// ForceStatus 直接改状态（管理 API / 测试）。
func ForceStatus(st Store, id, status string) error {
	rec, err := st.Get(id)
	if err != nil {
		return err
	}
	rec.Status = status
	return st.Update(rec)
}
