package store

import "github.com/LYH2263/go-outboxrelay/internal/codec"

// CloneRecord 深拷贝记录。
func CloneRecord(r Record) Record {
	return r
}

// CloneRecords 深拷贝切片。
func CloneRecords(in []Record) []Record {
	if in == nil {
		return nil
	}
	out := make([]Record, len(in))
	for i := range in {
		out[i] = CloneRecord(in[i])
	}
	return out
}
