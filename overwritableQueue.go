package utils

import "sync"

type OverWritableQueue[T any] struct {
	guard sync.RWMutex

	limit int
	queue []T
}

func NewOverWritableQueue[T any](limit int, preallocate bool) *OverWritableQueue[T] {
	var queue = &OverWritableQueue[T]{
		limit: limit,
	}
	if preallocate {
		queue.queue = make([]T, 0, limit)
	}
	return queue
}

func (o *OverWritableQueue[T]) Push(value T) bool {
	o.guard.Lock()
	defer o.guard.Unlock()

	if o.limit == 0 {
		return false
	}

	o.queue = append(o.queue, value)
	if len(o.queue) > o.limit {
		o.queue = o.queue[1:]
	}
	return true
}

// PushIfLenLessThanCap is an atomic version of the expression: if o.Len() < o.Cap() {o.Push(value)}
func (o *OverWritableQueue[T]) PushIfLenLessThanCap(value T) bool {
	o.guard.Lock()
	defer o.guard.Unlock()

	if o.limit == 0 {
		return false
	}

	if len(o.queue) < o.limit {
		o.queue = append(o.queue, value)
		return true
	} else {
		return false
	}
}

func (o *OverWritableQueue[T]) Filter(f func(v T) (keep bool)) {
	o.guard.Lock()
	defer o.guard.Unlock()

	var i = 0
	for _, v := range o.queue {
		if f(v) {
			o.queue[i] = v
			i++
		}
	}
	o.queue = o.queue[:i]
}

// Get returns the last item in the queue: [1, 2, 3] -> 3
func (o *OverWritableQueue[T]) Get() (val T, success bool) {
	o.guard.RLock()
	defer o.guard.RUnlock()

	if len(o.queue) != 0 {
		val, success = o.queue[len(o.queue)-1], true
	}
	return
}

// GetLeft returns the first item in the queue: [1, 2, 3] -> 1
func (o *OverWritableQueue[T]) GetLeft() (val T, success bool) {
	o.guard.RLock()
	defer o.guard.RUnlock()

	if len(o.queue) != 0 {
		val, success = o.queue[0], true
	}
	return
}

// Pull gets the last item in the queue and shifts it by 1 item: [1, 2, 3] -> item=3, queue=[1, 2]
func (o *OverWritableQueue[T]) Pull() (val T, success bool) {
	o.guard.Lock()
	defer o.guard.Unlock()

	if len(o.queue) != 0 {
		val, success = o.queue[len(o.queue)-1], true
		o.queue = o.queue[:len(o.queue)-1]
	}
	return
}

func (o *OverWritableQueue[T]) PullAndClear(deallocate bool) (val T, success bool) {
	o.guard.Lock()
	defer o.guard.Unlock()

	if len(o.queue) != 0 {
		val, success = o.queue[len(o.queue)-1], true
		o.clear(true, deallocate)
	}
	return
}

// ShiftLeft shifts the queue to the left by 1 element: [1, 2, 3] -> [2, 3]
func (o *OverWritableQueue[T]) ShiftLeft() bool {
	o.guard.Lock()
	defer o.guard.Unlock()

	if len(o.queue) != 0 {
		o.queue = o.queue[1:]
		return true
	}
	return false
}

func (o *OverWritableQueue[T]) Clear(deallocate bool) {
	o.clear(false, deallocate)
}

func (o *OverWritableQueue[T]) Len() int {
	o.guard.RLock()
	defer o.guard.RUnlock()
	return len(o.queue)
}

func (o *OverWritableQueue[T]) Cap() int {
	o.guard.RLock()
	defer o.guard.RUnlock()
	return o.limit
}

// Close clears internal data and prevents any writes
func (o *OverWritableQueue[T]) Close() {
	o.guard.Lock()
	defer o.guard.Unlock()

	o.limit = 0
	o.queue = nil
}

func (o *OverWritableQueue[T]) clear(externalLock bool, deallocate bool) {
	if !externalLock {
		o.guard.Lock()
		defer o.guard.Unlock()
	}

	if o.limit == 0 {
		return
	}

	if deallocate {
		o.queue = nil
	} else {
		o.queue = o.queue[:0]
	}
}
