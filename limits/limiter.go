package limits

import "sync"

type Limiter interface {
	Limit()
	Release()
}

type LimitWaiter interface {
	Limiter

	Wait()
}

type limiter struct {
	ch chan struct{}
	wg sync.WaitGroup
}

func NewLimiter(concurrency int) LimitWaiter {
	return &limiter{
		ch: make(chan struct{}, concurrency),
	}
}

func (lim *limiter) Limit() {
	lim.ch <- struct{}{}

	lim.wg.Add(1)
}

func (lim *limiter) Release() {
	lim.wg.Done()

	<-lim.ch
}

func (lim *limiter) Wait() {
	lim.wg.Wait()
}
