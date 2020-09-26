package utils

import "sync"

type SafeCounter struct {
	Value int
	s     sync.Mutex
}

func (c *SafeCounter) Increase() {
	c.s.Lock()
	c.Value++
	c.s.Unlock()
}

func (c *SafeCounter) Decrease() {
	c.s.Lock()
	c.Value--
	c.s.Unlock()
}

type SafeCounterLimited struct {
	MaxValue int
	DefValue int
	Value    int
	s        sync.Mutex
}

func (c *SafeCounterLimited) Increase() {
	c.s.Lock()
	if c.Value == c.MaxValue {
		c.Value = c.DefValue
	} else {
		c.Value++
	}
	c.s.Unlock()
}

func (c *SafeCounterLimited) Decrease() {
	c.s.Lock()
	c.Value--
	c.s.Unlock()
}
