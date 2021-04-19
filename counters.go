package utils

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"
)

func init() {
	a := make([]byte, 8)
	if _, e := cryptoRand.Read(a); e != nil {
		panic(e)
	}

	seed := binary.BigEndian.Uint64(a)
	rand.Seed(int64(seed))
}

type SafeCounter struct {
	sync.RWMutex

	waiters struct {
		sync.RWMutex
		m map[int64]*chan int
	}
	Value int
}

func NewSafeCounter() SafeCounter {
	return SafeCounter{
		RWMutex: sync.RWMutex{},
		waiters: struct {
			sync.RWMutex
			m map[int64]*chan int
		}{m: map[int64]*chan int{}},
		Value: 0,
	}
}

func (c *SafeCounter) Increase() {
	c.Lock()
	defer c.Unlock()
	c.Value++
	c.notify(c.Value)
}

func (c *SafeCounter) Decrease() {
	c.Lock()
	defer c.Unlock()
	c.Value--
	c.notify(c.Value)
}

func (c *SafeCounter) Get() int {
	c.RLock()
	defer c.RUnlock()
	return c.Value
}

func (c *SafeCounter) notify(num int) {
	c.waiters.RLock()
	for _, waiter := range c.waiters.m {
		*waiter <- num
	}
	c.waiters.RUnlock()
}

// param t: 0 - will return if value less or equals i, 1 - if value equals i, 2 - if value greater or equals i
type waitBehaviour int

const (
	WaitBehaviourLessOrEquals    waitBehaviour = 0
	WaitBehaviourEquals          waitBehaviour = 1
	WaitBehaviourGreaterOrEquals waitBehaviour = 2
)

func (c *SafeCounter) Wait(i int, behaviour waitBehaviour) {
	waiterKey := int64(rand.Uint64())
	ch := make(chan int, 1)
	ch <- c.Get()

	c.waiters.Lock()
	c.waiters.m[waiterKey] = &ch
	c.waiters.Unlock()

f:
	for {
		v := <-ch
		switch behaviour {
		case WaitBehaviourEquals:
			if v == i {
				break f
			}
		case WaitBehaviourLessOrEquals:
			if v <= i {
				break f
			}
		case WaitBehaviourGreaterOrEquals:
			if v >= i {
				break f
			}
		}
	}

	c.waiters.Lock()
	delete(c.waiters.m, waiterKey)
	c.waiters.Unlock()
}

type SafeCounterLimited struct {
	MaxValue int
	DefValue int
	Value    int
	s        sync.RWMutex
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

func (c *SafeCounterLimited) Get() int {
	c.s.RLock()
	defer c.s.RUnlock()
	return c.Value
}

//F
