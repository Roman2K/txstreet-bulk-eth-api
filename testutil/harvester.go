package testutil

type Harvester[T any] struct {
	ch      chan T
	results []T
}

func NewHarvester[T any]() *Harvester[T] {
	return &Harvester[T]{
		ch: make(chan T),
	}
}

func (h *Harvester[T]) Chan() chan<- T {
	return h.ch
}

func (h *Harvester[T]) Harvest() []T {
	if h.results != nil {
		return h.results
	}

	h.results = make([]T, 0)
	for elem := range h.ch {
		h.results = append(h.results, elem)
	}

	return h.results
}
