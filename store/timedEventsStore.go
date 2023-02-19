package store

import (
	"context"
	"sync"
	"time"
)

// TimedEventsStore stores events for a specified amount of time
type TimedEventsStore[EventT any] struct {
	s        sync.Mutex
	Events   map[int64]EventWrapper[EventT] // key - event expiration in unix nanos
	lifetime int64                          // nanoseconds
}

func NewTimedEventsStore[T any](lifetime time.Duration) *TimedEventsStore[T] {
	return &TimedEventsStore[T]{
		Events:   map[int64]EventWrapper[T]{},
		lifetime: lifetime.Nanoseconds(),
	}
}

func (t *TimedEventsStore[EventT]) add(now int64, event EventT) {
	t.clean(now)
	t.Events[now+t.lifetime] = EventWrapper[EventT]{EventTime: now, Event: event}
}

func (t *TimedEventsStore[EventT]) clean(now int64) {
	for expiration := range t.Events {
		if now > expiration {
			delete(t.Events, now)
		}
	}
}

// Tools

// Run cleans store in an indefinite loop. Call is not mandatory. Ticker will be stopped when ctx is done.
func (t *TimedEventsStore[EventT]) Run(ctx context.Context, ticker *time.Ticker) {
	defer ticker.Stop()
loop:
	for now := range ticker.C {
		select {
		case <-ctx.Done():
			break loop
		default:
			t.s.Lock()
			t.clean(now.UnixMilli())
			t.s.Unlock()
		}
	}
}

// Getters && setters

func (t *TimedEventsStore[EventT]) Count(externalLock bool) (count int) {
	if !externalLock {
		t.s.Lock()
		defer t.s.Unlock()
	}

	t.clean(time.Now().UnixNano())
	return len(t.Events)
}

func (t *TimedEventsStore[EventT]) Add(externalLock bool, event EventT) {
	if !externalLock {
		t.s.Lock()
		defer t.s.Unlock()
	}

	t.add(time.Now().UnixNano(), event)
}

// Impl. of sync.RWLocker

func (t *TimedEventsStore[EventT]) Lock() {
	t.s.Lock()
}
func (t *TimedEventsStore[EventT]) Unlock() {
	t.s.Unlock()
}
func (t *TimedEventsStore[EventT]) RLock() {
	t.s.Lock()
}
func (t *TimedEventsStore[EventT]) RUnlock() {
	t.s.Unlock()
}
