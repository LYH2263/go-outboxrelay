package store

import "github.com/LYH2263/go-outboxrelay/internal/codec"

// CloneRecord 深拷贝记录。
func CloneRecord(r Record) Record {
	out := r
	out.Payload = codec.CloneBytes(r.Payload)
	out.Headers = codec.CloneHeaders(r.Headers)
	return out
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
