package store

// Counts 各状态计数。
type Counts struct {
	Pending int
	Sending int
	Done    int
	Dead    int
	Total   int
}

// Summarize 汇总。
func Summarize(st Store) (Counts, error) {
	all, err := st.All()
	if err != nil {
		return Counts{}, err
	}
	var c Counts
	for _, r := range all {
		c.Total++
		switch r.Status {
		case "pending":
			c.Pending++
		case "sending":
			c.Sending++
		case "done":
			c.Done++
		case "dead":
			c.Dead++
		}
	}
	return c, nil
}
