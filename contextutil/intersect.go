package contextutil

import (
	"context"
	"time"
)

func Intersect(ctx, ctx2 context.Context) (context.Context, context.CancelFunc) {
	var cancel context.CancelFunc

	if deadline, ok := ctx2.Deadline(); ok {
		ctx, cancel = context.WithDeadline(ctx, deadline)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	go func() {
		select {
		case <-ctx.Done():
		case <-ctx2.Done():
			cancel()
		}
	}()

	return valuesUnionContext{ctx, ctx2}, cancel
}

type valuesUnionContext struct {
	ctx  context.Context
	ctx2 context.Context
}

func (vuc valuesUnionContext) Deadline() (time.Time, bool) { return vuc.ctx.Deadline() }
func (vuc valuesUnionContext) Done() <-chan struct{}       { return vuc.ctx.Done() }
func (vuc valuesUnionContext) Err() error                  { return vuc.ctx.Err() }

func (vuc valuesUnionContext) Value(key any) any {
	if value := vuc.ctx.Value(key); value != nil {
		return value
	}

	return vuc.ctx2.Value(key)
}
