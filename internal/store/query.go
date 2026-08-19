package store

import "strings"

// Filter 查询条件。
type Filter struct {
	Status  string
	Topic   string
	Limit   int
	Offset  int
}

// ApplyFilter 在已拷贝切片上过滤。
func ApplyFilter(recs []Record, f Filter) []Record {
	var out []Record
	skipped := 0
	for _, r := range recs {
		if f.Status != "" && r.Status != f.Status {
			continue
		}
		if f.Topic != "" && !topicMatch(r.Topic, f.Topic) {
			continue
		}
		if f.Offset > 0 && skipped < f.Offset {
			skipped++
			continue
		}
		out = append(out, r)
		if f.Limit > 0 && len(out) >= f.Limit {
			break
		}
	}
	return out
}

func topicMatch(topic, pattern string) bool {
	if pattern == "" || pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, ".*") {
		prefix := strings.TrimSuffix(pattern, ".*")
		return strings.HasPrefix(topic, prefix+".") || topic == prefix
	}
	return topic == pattern
}
