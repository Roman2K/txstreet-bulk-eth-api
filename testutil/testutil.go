package testutil

import (
	"context"
	"time"
)

func RunWithin(maxDuration time.Duration, fn func()) bool {
	success := false
	done := make(chan struct{})

	ctx, cancel := context.WithTimeout(context.Background(), maxDuration)
	defer cancel()

	go func() {
		fn()
		close(done)
	}()

	select {
	case <-done:
		success = true
	case <-ctx.Done():
	}

	return success
}
