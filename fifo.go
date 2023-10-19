package utils

import "sync"

type Fifo[T any] struct {
	s     sync.Mutex
	items []*T
}

func (f *Fifo[T]) Pull() (T, bool) {
	f.s.Lock()
	defer f.s.Unlock()

	if len(f.items) == 0 {
		var v T
		return v, false
	}
	v := f.items[0]
	f.items = f.items[1:]
	return *v, true
}

func (f *Fifo[T]) HasAny() bool {
	f.s.Lock()
	defer f.s.Unlock()

	return len(f.items) != 0
}

func (f *Fifo[T]) Push(item T) {
	f.s.Lock()
	defer f.s.Unlock()
	f.items = append(f.items, &item)
}
