package store

import (
	"context"
	"sync"
	"time"
)

// TimedMultiEventsStoreL3 stores multiple events of different subtypes and sub-subtypes for a specified amount of time
type TimedMultiEventsStoreL3[KL3, KL2 comparable, EventT any] struct {
	s        sync.Mutex
	Events   map[KL3]map[KL2]map[int64]EventWrapper[EventT] // event sub-subtype (L3) -> event subtype (L2) -> event expiration in unix nanos (l1) -> event
	lifetime int64                                          // nanoseconds

	DisableChecksOnInserts bool
	//DisableChecksOnRetrieves bool
}

func NewTimedMultiEventsStoreL3[KL3, KL2 comparable, T any](lifetime time.Duration) *TimedMultiEventsStoreL3[KL3, KL2, T] {
	return &TimedMultiEventsStoreL3[KL3, KL2, T]{
		Events:   map[KL3]map[KL2]map[int64]EventWrapper[T]{},
		lifetime: lifetime.Nanoseconds(),
	}
}

func (t *TimedMultiEventsStoreL3[KL3, KL2, EventT]) count(now int64, subkey KL3, key KL2, deleteMapWhenEmpty bool) (count int) {
	//if !t.DisableChecksOnRetrieves {
	//	t.clean(now, subkey, key, deleteMapWhenEmpty)
	//}
	if v, h := t.Events[subkey]; h {
		if v, h := v[key]; h {
			return len(v)
		}
	}
	return 0
}

// has returns whether the keys exist in the store. Doesn't check for a 2lvl map length != 0. In you need to, you need count() != 0
func (t *TimedMultiEventsStoreL3[KL3, KL2, EventT]) has(now int64, subkey KL3, key KL2, deleteMapWhenEmpty bool) (has bool) {
	//if !t.DisableChecksOnRetrieves {
	//	t.clean(now, subkey, key, deleteMapWhenEmpty)
	//}
	if v, h := t.Events[subkey]; h {
		_, h := v[key]
		return h
	} else {
		return false
	}
}

func (t *TimedMultiEventsStoreL3[KL3, KL2, EventT]) add(now int64, subkey KL3, key KL2, event EventT) {
	//if !t.DisableChecksOnRetrieves {
	//	t.clean(now, subkey, key, false)
	//}
	vl3, h := t.Events[subkey]
	if !h {
		vl3 = map[KL2]map[int64]EventWrapper[EventT]{}
		t.Events[subkey] = vl3
	}
	vl2, h := vl3[key]
	if !h {
		vl2 = map[int64]EventWrapper[EventT]{}
		vl3[key] = vl2
	}
	vl2[now+t.lifetime] = EventWrapper[EventT]{EventTime: now, Event: event}
}

// clean cleans only l1 map (and deletes l2 (key) map from l3 (subkey) map if l1 (time) map is empty and deleteMap is set to true)
// calling clean with non-existing subkey and key is safe
func (t *TimedMultiEventsStoreL3[KL3, KL2, EventT]) clean(now int64, subkey KL3, key KL2, deleteMap bool) {
	if mL3, h := t.Events[subkey]; h {
		if mL2, h := mL3[key]; h {
			var valid bool
			for expiration := range mL2 {
				if now > expiration {
					delete(mL2, now)
				} else if !valid {
					valid = true
				}
			}
			if deleteMap && !valid {
				delete(mL3, key)
			}
		}
	}
}

// cleanL3 cleans l2, l1 maps of the l3 subkey. Subkey must exist in the l3 map.
func (t *TimedMultiEventsStoreL3[KL3, KL2, EventT]) cleanL3(now int64, subkey KL3, deleteMap bool) {
	var mL3 = t.Events[subkey]
	for key, mL2 := range mL3 {
		var valid bool
		for expiration := range mL2 {
			if now > expiration {
				delete(mL2, now)
			} else if !valid {
				valid = true
			}
		}
		if deleteMap && !valid {
			delete(mL3, key)
		}
	}
}

// Tools

// Run cleans store in an indefinite loop. Call is not mandatory. Ticker will be stopped when ctx is done.
// deleteMapWhenEmpty specifies whether the sub-map should be deleted when it's empty.
func (t *TimedMultiEventsStoreL3[KL3, KL2, EventT]) Run(ctx context.Context, ticker *time.Ticker, deleteMapWhenEmpty bool) {
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
				t.cleanL3(now, k, deleteMapWhenEmpty)
			}
			t.s.Unlock()
		}
	}
}

// Getters && setters

// Count
// deleteMapWhenEmpty specifies whether the sub-map should be deleted when it's empty.
func (t *TimedMultiEventsStoreL3[KL3, KL2, EventT]) Count(externalLock bool, subKey KL3, key KL2, deleteMapWhenEmpty bool) (count int) {
	if !externalLock {
		t.s.Lock()
		defer t.s.Unlock()
	}

	return t.count(time.Now().UnixNano(), subKey, key, deleteMapWhenEmpty)
}

// Has returns whether the key exists in the store regardless of its content length.
// deleteMapWhenEmpty specifies whether the sub-map should be deleted when it's empty.
func (t *TimedMultiEventsStoreL3[KL3, KL2, EventT]) Has(externalLock bool, subKey KL3, key KL2, deleteMapWhenEmpty bool) (has bool) {
	if !externalLock {
		t.s.Lock()
		defer t.s.Unlock()
	}

	return t.has(time.Now().UnixNano(), subKey, key, deleteMapWhenEmpty)
}

// HasF is the same as Has, but the caller has full control (locking is up to the caller too)
func (t *TimedMultiEventsStoreL3[KL3, KL2, EventT]) HasF(now int64, subKey KL3, key KL2, deleteMapWhenEmpty bool) bool {
	return t.has(now, subKey, key, deleteMapWhenEmpty)
}

func (t *TimedMultiEventsStoreL3[KL3, KL2, EventT]) Add(externalLock bool, subKey KL3, key KL2, event EventT) {
	if !externalLock {
		t.s.Lock()
		defer t.s.Unlock()
	}

	t.add(time.Now().UnixNano(), subKey, key, event)
}

// Impl. of sync.RWLocker

func (t *TimedMultiEventsStoreL3[KL3, KL2, EventT]) Lock() {
	t.s.Lock()
}
func (t *TimedMultiEventsStoreL3[KL3, KL2, EventT]) Unlock() {
	t.s.Unlock()
}
func (t *TimedMultiEventsStoreL3[KL3, KL2, EventT]) RLock() {
	t.s.Lock()
}
func (t *TimedMultiEventsStoreL3[KL3, KL2, EventT]) RUnlock() {
	t.s.Unlock()
}
