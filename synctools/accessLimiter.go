package synctools

import (
	"sync"
)

type AccessLimiter struct {
	cond    sync.Cond
	running int

	maxSimultaneouslyRunning int
}

func NewAccessLimiter(maxSimultaneouslyRunning int) *AccessLimiter {
	if maxSimultaneouslyRunning == 0 {
		panic("maxSimultaneouslyRunning must be > 0")
	}

	return &AccessLimiter{cond: sync.Cond{L: &sync.Mutex{}}, maxSimultaneouslyRunning: maxSimultaneouslyRunning}
}

func (a *AccessLimiter) Queue(f func()) {
	a.cond.L.Lock()
	for a.running >= a.maxSimultaneouslyRunning {
		a.cond.Wait()
	}
	a.running++
	a.cond.L.Unlock()

	f()
	a.cond.L.Lock()
	a.running--
	a.cond.Signal()
	a.cond.L.Unlock()
}
