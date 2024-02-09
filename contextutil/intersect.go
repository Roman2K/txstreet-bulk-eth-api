package contextutil

import (
	"context"
	"time"
)

func Intersect(ctx, ctx2 context.Context) (context.Context, context.CancelFunc) {
	var cancelDeadline context.CancelFunc

	if deadline, ok := ctx2.Deadline(); ok {
		ctx, cancelDeadline = context.WithDeadline(ctx, deadline)
	}

	var cancelCause context.CancelCauseFunc
	ctx, cancelCause = context.WithCancelCause(ctx)

	go func() {
		if cancelDeadline != nil {
			defer cancelDeadline()
		}

		select {
		case <-ctx.Done():
		case <-ctx2.Done():
			cancelCause(context.Cause(ctx2))
		}
	}()

	return valuesUnionContext{ctx, ctx2}, func() { cancelCause(nil) }
}

type valuesUnionContext struct {
	ctx  context.Context
	ctx2 context.Context
}

func (vuc valuesUnionContext) Deadline() (time.Time, bool) { return vuc.ctx.Deadline() }
func (vuc valuesUnionContext) Done() <-chan struct{}       { return vuc.ctx.Done() }

func (vuc valuesUnionContext) Err() error {
	err := vuc.ctx.Err()
	ctx2Err := vuc.ctx2.Err()

	if err == context.Canceled && ctx2Err == context.DeadlineExceeded {
		return ctx2Err
	}

	return err
}

func (vuc valuesUnionContext) Value(key any) any {
	if value := vuc.ctx.Value(key); value != nil {
		return value
	}

	return vuc.ctx2.Value(key)
}
