package outboxrelay

// Close 标记关闭：先 Flush 持久化，再关闭并释放 store，最后清空投递客户端。
// 之后 Append / RelayOnce / ListPending 返回 ErrClosed，不会解引用 nil store。
func (o *Outbox) Close() error {
	if o == nil {
		return ErrNilOutbox
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return nil
	}
	o.closed = true

	var first error
	if o.st != nil {
		if err := o.st.Close(); err != nil && first == nil {
			first = err
		}
		if err := o.st.Flush(); err != nil && first == nil {
			first = err
		}
	}
	o.st = nil
	if o.client != nil {
		o.client.CloseIdle()
	}
	o.client = nil
	o.deliv = nil
	return first
}

// Closed 报告是否已关闭。
func (o *Outbox) Closed() bool {
	if o == nil {
		return true
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.closed
}
