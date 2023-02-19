package store

import (
	"context"
	"sync"
	"time"
)

// TimedMultiEventsStoreL2 stores multiple events of different subtypes for a specified amount of time
type TimedMultiEventsStoreL2[K comparable, EventT any] struct {
	s        sync.Mutex
	Events   map[K]map[int64]EventWrapper[EventT] // event subtype (L2) -> event expiration in unix nanos -> event
	lifetime int64                                // nanoseconds

	DisableChecksOnInserts   bool
	DisableChecksOnRetrieves bool
}

func NewTimedMultiEventsStore[K comparable, T any](lifetime time.Duration) *TimedMultiEventsStoreL2[K, T] {
	return &TimedMultiEventsStoreL2[K, T]{
		Events:   map[K]map[int64]EventWrapper[T]{},
		lifetime: lifetime.Nanoseconds(),
	}
}

func (t *TimedMultiEventsStoreL2[K, EventT]) count(key K, deleteMapWhenEmpty bool) (count int) {
	if !t.DisableChecksOnRetrieves {
		t.clean(time.Now().UnixNano(), key, deleteMapWhenEmpty)
	}
	if v, h := t.Events[key]; h {
		return len(v)
	} else {
		return 0
	}
}

// has returns whether the key exists in the store. Doesn't check for a 2lvl map length != 0. In you need to, you need count() != 0
func (t *TimedMultiEventsStoreL2[K, EventT]) has(key K, deleteMapWhenEmpty bool) (has bool) {
	if !t.DisableChecksOnRetrieves {
		t.clean(time.Now().UnixNano(), key, deleteMapWhenEmpty)
	}
	_, h := t.Events[key]
	return h
}

func (t *TimedMultiEventsStoreL2[K, EventT]) add(now int64, key K, event EventT) {
	if !t.DisableChecksOnRetrieves {
		t.clean(now, key, false)
	}
	v, h := t.Events[key]
	if !h {
		v = map[int64]EventWrapper[EventT]{}
		t.Events[key] = v
	}
	v[now+t.lifetime] = EventWrapper[EventT]{EventTime: now, Event: event}
}

func (t *TimedMultiEventsStoreL2[K, EventT]) clean(now int64, key K, deleteMap bool) {
	var valid bool
	for expiration := range t.Events[key] {
		if now > expiration {
			delete(t.Events[key], now)
		} else if !valid {
			valid = true
		}
	}
	if deleteMap && !valid {
		delete(t.Events, key)
	}
}

// Tools

// Run cleans store in an indefinite loop. Call is not mandatory. Ticker will be stopped when ctx is done.
// deleteMapWhenEmpty specifies whether the sub-map should be deleted when it's empty.
func (t *TimedMultiEventsStoreL2[K, EventT]) Run(ctx context.Context, ticker *time.Ticker, deleteMapWhenEmpty bool) {
	defer ticker.Stop()
loop:
	for now := range ticker.C {
		select {
		case <-ctx.Done():
			break loop
		default:
			var now = now.UnixMilli()
			t.s.Lock()
			for k := range t.Events {
				t.clean(now, k, deleteMapWhenEmpty)
			}
			t.s.Unlock()
		}
	}
}

// Getters && setters

// Count
// deleteMapWhenEmpty specifies whether the sub-map should be deleted when it's empty.
func (t *TimedMultiEventsStoreL2[K, EventT]) Count(externalLock bool, key K, deleteMapWhenEmpty bool) (count int) {
	if !externalLock {
		t.s.Lock()
		defer t.s.Unlock()
	}

	return t.count(key, deleteMapWhenEmpty)
}

// Has returns whether the key exists in the store regardless of its content length.
// deleteMapWhenEmpty specifies whether the sub-map should be deleted when it's empty.
func (t *TimedMultiEventsStoreL2[K, EventT]) Has(externalLock bool, key K, deleteMapWhenEmpty bool) (has bool) {
	if !externalLock {
		t.s.Lock()
		defer t.s.Unlock()
	}

	return t.has(key, deleteMapWhenEmpty)
}

func (t *TimedMultiEventsStoreL2[K, EventT]) Add(externalLock bool, key K, event EventT) {
	if !externalLock {
		t.s.Lock()
		defer t.s.Unlock()
	}

	t.add(time.Now().UnixNano(), key, event)
}

// Impl. of sync.RWLocker

func (t *TimedMultiEventsStoreL2[K, EventT]) Lock() {
	t.s.Lock()
}
func (t *TimedMultiEventsStoreL2[K, EventT]) Unlock() {
	t.s.Unlock()
}
func (t *TimedMultiEventsStoreL2[K, EventT]) RLock() {
	t.s.Lock()
}
func (t *TimedMultiEventsStoreL2[K, EventT]) RUnlock() {
	t.s.Unlock()
}
