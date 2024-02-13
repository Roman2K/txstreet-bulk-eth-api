package collection

import (
	"errors"
	"sync"

	"github.com/Roman2K/txstreet-bulk-eth-api/limits"
)

type Collector[InT, OutT any] interface {
	CollectAll([]InT) ([]OutT, error)
}

type CollectFunc[InT, OutT any] func(InT) (OutT, error)

type FuncCollector[InT, OutT any] CollectFunc[InT, OutT]

func (fn FuncCollector[InT, OutT]) CollectAll(elems []InT) ([]OutT, error) {
	results := make([]OutT, len(elems))
	errs := make([]error, len(elems))

	var wg sync.WaitGroup
	wg.Add(len(elems))

	for index := range elems {
		go func(index int) {
			defer wg.Done()

			results[index], errs[index] = fn(elems[index])
		}(index)
	}

	wg.Wait()

	return results, errors.Join(errs...)
}

type LimitCollector[InT, OutT any] struct {
	Limiter     limits.Limiter
	CollectFunc CollectFunc[InT, OutT]
}

func (c LimitCollector[InT, OutT]) CollectAll(elems []InT) ([]OutT, error) {
	collect := func(elem InT) (OutT, error) {
		c.Limiter.Limit()
		defer c.Limiter.Release()

		return c.CollectFunc(elem)
	}

	return FuncCollector[InT, OutT](collect).CollectAll(elems)
}

var (
	_ Collector[int, bool] = FuncCollector[int, bool](func(int) (bool, error) { return false, nil })
	_ Collector[int, bool] = LimitCollector[int, bool]{}
)
