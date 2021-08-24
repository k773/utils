package maps

import (
	"sync"
	"time"
)

type LockExpireMap struct {
	s sync.Mutex
	m map[string]*sync.Mutex
	t map[time.Time][]string
}

func NewLockExpireMap() *LockExpireMap {
	a := &LockExpireMap{
		m: map[string]*sync.Mutex{},
		t: map[time.Time][]string{},
	}
	go a.StartAutoRemoval()
	return a
}

func (l *LockExpireMap) StartAutoRemoval() {
	for t := range time.NewTicker(30 * time.Second).C {
		l.s.Lock()
		k := t.Truncate(time.Minute)
		for _, k := range l.t[k] {
			delete(l.m, k)
		}
		delete(l.t, k)
		l.s.Unlock()
	}
}

// Get : Lifetime must be greater than minute
func (l *LockExpireMap) Get(k string, lifetime time.Duration) (v *sync.Mutex) {
	l.s.Lock()
	defer l.s.Unlock()
	var h bool
	if v, h = l.m[k]; !h {
		v = &sync.Mutex{}
		l.m[k] = v
		t := time.Now().Add(lifetime).Truncate(time.Minute)
		l.t[t] = append(l.t[t], k)
	}
	return
}
