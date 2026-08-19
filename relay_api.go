package outboxrelay

import (
	"context"
	"errors"
	"time"

	"github.com/LYH2263/go-outboxrelay/internal/backoff"
	"github.com/LYH2263/go-outboxrelay/internal/deliver"
	"github.com/LYH2263/go-outboxrelay/internal/relay"
	"github.com/LYH2263/go-outboxrelay/internal/status"
	"github.com/LYH2263/go-outboxrelay/internal/store"
)

// RelayOnce 拉取一条到期 pending 事件并投递一次。
func (o *Outbox) RelayOnce(ctx context.Context) (RelayResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	o.mu.Lock()
	if err := o.checkOpenLocked(); err != nil {
		o.mu.Unlock()
		return RelayResult{}, err
	}
	if o.st == nil {
		o.mu.Unlock()
		return RelayResult{}, ErrClosed
	}
	st := o.st
	deliv := o.effectiveDelivererLocked()
	pol := o.policy
	clk := o.clk
	o.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return RelayResult{}, mapCtxErr(err)
	}

	rec, ok, err := relay.ClaimNext(st, clk.Now())
	if err != nil {
		return RelayResult{}, err
	}
	if !ok {
		return RelayResult{}, nil
	}

	res, err := o.deliverOne(ctx, deliv, pol, rec)
	o.mu.Lock()
	o.relayed++
	o.lastRelay = clk.Now()
	if res.OK {
		o.succeeded++
	} else if res.Err != "" {
		o.failed++
	}
	o.mu.Unlock()
	return res, err
}

// RelayContext 在 ctx 取消前最多投递 max 条（max<=0 表示不限，但仍尊重取消）。
func (o *Outbox) RelayContext(ctx context.Context, max int) ([]RelayResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var out []RelayResult
	for {
		if err := ctx.Err(); err != nil {
			if len(out) == 0 {
				return out, mapCtxErr(err)
			}
			return out, mapCtxErr(err)
		}
		if max > 0 && len(out) >= max {
			return out, nil
		}
		res, err := o.RelayOnce(ctx)
		if err != nil {
			if errors.Is(err, ErrClosed) || errors.Is(err, ErrCanceled) || errors.Is(err, context.Canceled) {
				return out, err
			}
			if res.EventID == "" && err != nil {
				return out, err
			}
		}
		if res.EventID == "" {
			return out, nil
		}
		out = append(out, res)
		// 失败后退避（尊重 ctx）
		if !res.OK && res.Attempts > 0 {
			o.mu.Lock()
			pol := o.policy
			o.mu.Unlock()
			d := backoff.Delay(pol, res.Attempts)
			if err := backoff.Wait(ctx, d); err != nil {
				return out, mapCtxErr(err)
			}
		}
	}
}

func (o *Outbox) deliverOne(ctx context.Context, deliv Deliverer, pol backoff.Policy, rec store.Record) (RelayResult, error) {
	start := time.Now()
	res := RelayResult{
		EventID:  rec.ID,
		Topic:    rec.Topic,
		Attempts: rec.Attempts + 1,
	}
	code, err := deliv.Deliver(ctx, rec.Topic, rec.Payload, rec.Headers, rec.TargetURL)
	res.Duration = time.Since(start)
	res.StatusCode = code
	res.FinishedAt = o.clk.Now()

	if err == nil && code >= 200 && code < 300 {
		res.OK = true
		if merr := status.MarkDone(o.st, rec.ID); merr != nil {
			res.OK = false
			res.Err = merr.Error()
			return res, merr
		}
		return res, nil
	}

	msg := ""
	if err != nil {
		msg = err.Error()
		res.Err = msg
	} else {
		msg = "non-2xx status"
		res.Err = msg
	}
	dead, merr := status.MarkFail(o.st, rec.ID, msg, pol, o.clk.Now())
	if merr != nil {
		return res, merr
	}
	if dead {
		o.mu.Lock()
		o.dead++
		o.mu.Unlock()
	}
	return res, nil
}

func (o *Outbox) effectiveDelivererLocked() Deliverer {
	if o.deliverer != nil {
		return o.deliverer
	}
	return &deliverAdapter{d: o.deliv}
}

type deliverAdapter struct {
	d *deliver.HTTPDeliverer
}

func (a *deliverAdapter) Deliver(ctx context.Context, topic string, payload []byte, headers map[string]string, url string) (int, error) {
	if a == nil || a.d == nil {
		return 0, ErrNilTransport
	}
	return a.d.Post(ctx, topic, payload, headers, url)
}

func mapCtxErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) {
		return ErrCanceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeout
	}
	return err
}
