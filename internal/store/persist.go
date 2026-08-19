package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type snapshot struct {
	Records []Record `json:"records"`
}

func (m *Memory) persistLocked() error {
	if m.path == "" {
		return nil
	}
	if m.flushErr != nil {
		err := m.flushErr
		m.flushErr = nil
		return err
	}
	recs := make([]Record, 0, len(m.order))
	for _, id := range m.order {
		if r := m.byID[id]; r != nil {
			recs = append(recs, CloneRecord(*r))
		}
	}
	data, err := json.MarshalIndent(snapshot{Records: recs}, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(m.path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	tmp := m.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, m.path)
}

func (m *Memory) loadLocked() error {
	data, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var snap snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return err
	}
	m.byID = make(map[string]*Record, len(snap.Records))
	m.order = m.order[:0]
	for _, r := range snap.Records {
		cp := CloneRecord(r)
		m.byID[cp.ID] = &cp
		m.order = append(m.order, cp.ID)
	}
	return nil
}

// SaveJSON 写出任意记录切片（供外部 sink）。
func SaveJSON(path string, recs []Record) error {
	data, err := json.MarshalIndent(snapshot{Records: CloneRecords(recs)}, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// LoadJSON 读入快照。
func LoadJSON(path string) ([]Record, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var snap snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	return CloneRecords(snap.Records), nil
}
