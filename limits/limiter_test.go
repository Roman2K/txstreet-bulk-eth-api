package limits

import (
	"sync"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"

	"github.com/Roman2K/txstreet-bulk-eth-api/testutil"
)

func TestLimitWaiter(t *testing.T) {
	limiter := NewLimiter(1)

	executed := false
	noTimeout := testutil.RunWithin(100*time.Millisecond, func() {
		limiter.Limit()
		limiter.Release()

		executed = true
	})

	assert.True(t, noTimeout)
	assert.True(t, executed)
}

func TestLimitWaiterHittingLimit(t *testing.T) {
	limiter := NewLimiter(1)
	results := testutil.NewHarvester[int]()

	noTimeout := testutil.RunWithin(70*time.Millisecond, func() {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			results.Chan() <- 1

			limiter.Limit()
			results.Chan() <- 2

			limiter.Limit()
			results.Chan() <- 4

			limiter.Release()
		}()

		go func() {
			defer wg.Done()

			time.Sleep(50 * time.Millisecond)

			results.Chan() <- 3

			limiter.Release()
		}()

		go func() {
			wg.Wait()
			close(results.Chan())
		}()

		results.Harvest()
	})

	assert.True(t, noTimeout)
	assert.Equal(t, []int{1, 2, 3, 4}, results.Harvest())
}

func TestLimitWaiterWait(t *testing.T) {
	limiter := NewLimiter(1)
	startTime := time.Now()

	go func() {
		limiter.Limit()
		limiter.Limit()
		limiter.Release()
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		limiter.Release()
	}()

	noTimeout := testutil.RunWithin(70*time.Millisecond, func() {
		limiter.Wait()
	})
	endTime := time.Now()

	assert.True(t, noTimeout)

	assert.WithinDuration(
		t,
		endTime,
		startTime.Add(50*time.Millisecond),
		5*time.Millisecond,
	)
}
