package store

import "github.com/LYH2263/go-outboxrelay/internal/codec"

// CloneRecord 深拷贝记录：值类型字段随结构体复制，引用类型字段
// （Payload 切片、Headers map）必须单独拷贝底层数据，否则调用方
// 就地改写返回值会写穿到 store 内存里的原始记录。
func CloneRecord(r Record) Record {
	r.Payload = codec.CloneBytes(r.Payload)
	r.Headers = codec.CloneHeaders(r.Headers)
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
