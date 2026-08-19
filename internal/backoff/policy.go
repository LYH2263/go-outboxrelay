package backoff

import "time"

// Policy 指数退避。
type Policy struct {
	Base   time.Duration
	Cap    time.Duration
	Factor float64
}

// Default 默认策略。
func Default() Policy {
	return Policy{Base: 20 * time.Millisecond, Cap: 2 * time.Second, Factor: 2}
}

func (p Policy) normalized() Policy {
	if p.Base <= 0 {
		p.Base = 20 * time.Millisecond
	}
	if p.Cap <= 0 {
		p.Cap = 2 * time.Second
	}
	if p.Factor < 1 {
		p.Factor = 2
	}
	return p
}
