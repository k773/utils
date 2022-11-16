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
		queue.queue = make([]T, limit, 0)
	}
	return queue
}

func (o *OverWritableQueue[T]) Add(value T) bool {
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

func (o *OverWritableQueue[T]) Get() (val T, success bool) {
	o.guard.RLock()
	defer o.guard.RUnlock()

	if len(o.queue) != 0 {
		val, success = o.queue[len(o.queue)-1], true
	}
	return
}

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

func (o *OverWritableQueue[T]) Clear(deallocate bool) {
	o.clear(false, deallocate)
}

func (o *OverWritableQueue[T]) clear(externalLock bool, deallocate bool) {
	if !externalLock {
		o.guard.Lock()
		defer o.guard.Unlock()
	}

	if deallocate {
		o.queue = nil
	} else {
		o.queue = o.queue[:0]
	}
}
