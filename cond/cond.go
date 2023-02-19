// Implementation of simple Waiter-Broadcast model

package cond

import (
	"context"
	"sync"
)

// Cond
// Warning: in the current implementation nothing prevents data from changing in between the moment we receive notification from the channel
// and the moment we obtain a data lock (c.L.Lock() in Wait()). Use sync/cond.
type Cond struct {
	L  sync.Locker
	ch chan struct{}
}

func NewCond(l sync.Locker) *Cond {
	return &Cond{L: l, ch: make(chan struct{})}
}

func (c *Cond) Wait(ctx context.Context) {
	var ch = c.ch
	c.L.Unlock()
	select {
	case <-ch:
	case <-ctx.Done():
	}
	c.L.Lock()
}

func (c *Cond) Broadcast() {
	var ch = make(chan struct{})
	c.L.Lock()
	var old = c.ch
	c.ch = ch
	close(old)
	c.L.Unlock()
}
