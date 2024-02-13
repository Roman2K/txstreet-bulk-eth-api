package contextutil

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Roman2K/txstreet-bulk-eth-api/testutil"
	assert "github.com/stretchr/testify/require"
)

func TestIntersectCancellingParent(t *testing.T) {
	parent, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	ctx2, cancel2 := context.WithCancelCause(context.Background())
	defer cancel2(nil)

	intersection, cancelIntersection := Intersect(parent, ctx2)
	defer cancelIntersection()

	assert.False(t, isDone(intersection))

	cancel1()

	assert.True(t, isDone(intersection))

	assert.Equal(t, context.Canceled, intersection.Err())
	assert.Equal(t, context.Canceled, context.Cause(intersection))
}

// })

func TestIntersectCancellingOtherWithCause(t *testing.T) {
	parent, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	ctx2, cancel2 := context.WithCancelCause(context.Background())
	defer cancel2(nil)

	intersection, cancelIntersection := Intersect(parent, ctx2)
	defer cancelIntersection()

	assert.False(t, isDone(intersection))

	someCause := errors.New("Some cause")
	cancel2(someCause)

	assert.True(t, isDone(intersection))

	assert.Equal(t, context.Canceled, intersection.Err())
	assert.Equal(t, someCause, context.Cause(intersection))
}

func TestIntersectWithOtherContextDeadlineCancellingOther(t *testing.T) {
	parent, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	intersection, cancelIntersection := Intersect(parent, ctx2)
	defer cancelIntersection()

	assert.False(t, isDone(intersection))

	cancel2()

	assert.True(t, isDone(intersection))

	assert.Equal(t, context.Canceled, intersection.Err())
	assert.Equal(t, context.Canceled, context.Cause(intersection))
}

func TestIntersectWithOtherContextDeadlineOtherTimesOut(t *testing.T) {
	parent, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()

	intersection, cancelIntersection := Intersect(parent, ctx2)
	defer cancelIntersection()

	assert.False(t, isDone(intersection))

	<-intersection.Done()

	assert.True(t, isDone(intersection))

	assert.Equal(t, context.DeadlineExceeded, intersection.Err())
	assert.Equal(t, context.DeadlineExceeded, context.Cause(intersection))
}

func TestIntersectWithOtherContextDeadlineParentTimesOut(t *testing.T) {
	parent, cancel1 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel1()

	ctx2, cancel2 := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel2()

	intersection, cancelIntersection := Intersect(parent, ctx2)
	defer cancelIntersection()

	assert.False(t, isDone(intersection))

	noTimeout := testutil.RunWithin(60*time.Millisecond, func() {
		<-intersection.Done()
	})

	assert.True(t, noTimeout)
	assert.True(t, isDone(intersection))

	assert.Equal(t, context.DeadlineExceeded, intersection.Err())
	assert.Equal(t, context.DeadlineExceeded, context.Cause(intersection))

	noTimeout = testutil.RunWithin(30*time.Millisecond, func() {
		<-ctx2.Done()
	})

	assert.True(t, noTimeout)

	assert.Equal(t, context.DeadlineExceeded, intersection.Err())
	assert.Equal(t, context.DeadlineExceeded, context.Cause(intersection))
}

func isDone(ctx context.Context) bool {
	timeoutCtx, cancel :=
		context.WithTimeout(context.Background(), 50*time.Millisecond)

	defer cancel()

	select {
	case <-timeoutCtx.Done():
	case <-ctx.Done():
		return true
	}
	return false
}

func TestIntersectTwoDeadlines(t *testing.T) {
	now := time.Now()

	parent, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	intersection, cancelIntersection := Intersect(parent, ctx2)
	defer cancelIntersection()

	deadline, ok := intersection.Deadline()
	assert.True(t, ok)

	assert.WithinDuration(t, now.Add(5*time.Second), deadline, 50*time.Millisecond)
}

func TestIntersectOnlyParentDeadline(t *testing.T) {
	now := time.Now()

	parent, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	ctx2 := context.Background()

	intersection, cancelIntersection := Intersect(parent, ctx2)
	defer cancelIntersection()

	deadline, ok := intersection.Deadline()
	assert.True(t, ok)

	assert.WithinDuration(t, now.Add(10*time.Second), deadline, 50*time.Millisecond)
}

func TestIntersectOnlyOtherContextDeadline(t *testing.T) {
	now := time.Now()

	parent := context.Background()

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	intersection, cancelIntersection := Intersect(parent, ctx2)
	defer cancelIntersection()

	deadline, ok := intersection.Deadline()
	assert.True(t, ok)

	assert.WithinDuration(t, now.Add(5*time.Second), deadline, 50*time.Millisecond)
}

func TestIntersectValues(t *testing.T) {
	parent := context.WithValue(context.Background(), "key1", "value1")
	ctx2 := context.WithValue(context.Background(), "key2", "value2")

	intersection, cancelIntersection := Intersect(parent, ctx2)
	defer cancelIntersection()

	assert.Nil(t, intersection.Value("xxx"))

	assert.Equal(t, "value1", intersection.Value("key1"))
	assert.Equal(t, "value2", intersection.Value("key2"))
}
