package synctools

import (
	"context"
	"sync"
)

type ChanMutexLazyInit struct {
	s sync.Mutex
	m *ChanMutex
}

func (c *ChanMutexLazyInit) init() {
	if c.m != nil {
		return
	}
	c.s.Lock()
	defer c.s.Unlock()
	if c.m == nil {
		c.m = NewChanMutex()
	}
}

func (c *ChanMutexLazyInit) TryLockCtx(ctx context.Context) error {
	c.init()
	return c.m.TryLockCtx(ctx)
}

func (c *ChanMutexLazyInit) TryLock(stop <-chan struct{}) bool {
	c.init()
	return c.m.TryLock(stop)
}

func (c *ChanMutexLazyInit) Unlock() {
	c.m.Unlock()
}

type ChanMutex struct {
	c chan struct{}
}

func NewChanMutex() *ChanMutex {
	return &ChanMutex{c: make(chan struct{}, 1)}
}

func (c *ChanMutex) TryLockCtx(ctx context.Context) error {
	if c.TryLock(ctx.Done()) {
		return nil
	}
	return ctx.Err()
}

func (c *ChanMutex) TryLock(stop <-chan struct{}) bool {
	select {
	case c.c <- struct{}{}:
		return true
	case <-stop:
		return false
	}
}

func (c *ChanMutex) Unlock() {
	<-c.c
}
