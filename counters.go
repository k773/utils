package utils

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"
	"time"
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
	notifier *sync.Cond
	Value    int
}

func NewSafeCounter() SafeCounter {
	return SafeCounter{
		RWMutex:  sync.RWMutex{},
		notifier: sync.NewCond(&sync.Mutex{}),
		Value:    0,
	}
}

func (c *SafeCounter) Increase() {
	c.Lock()
	c.notifier.L.Lock()
	defer c.Unlock()
	defer c.notifier.L.Unlock()
	c.Value++
	c.notify()
}

func (c *SafeCounter) Decrease() {
	c.Lock()
	c.notifier.L.Lock()
	defer c.Unlock()
	defer c.notifier.L.Unlock()
	c.Value--
	c.notify()
}

func (c *SafeCounter) Get() int {
	c.RLock()
	defer c.RUnlock()
	return c.Value
}

func (c *SafeCounter) notify() {
	c.notifier.Broadcast()
}

// param t: 0 - will return if value less or equals i, 1 - if value equals i, 2 - if value greater or equals i
type waitBehaviour int

const (
	WaitBehaviourLessOrEquals    waitBehaviour = 0
	WaitBehaviourEquals          waitBehaviour = 1
	WaitBehaviourGreaterOrEquals waitBehaviour = 2
)

func (c *SafeCounter) Wait(i int, behaviour waitBehaviour) {
	//c.notifier.L.Lock()
	//defer c.notifier.L.Unlock()
f:
	for {
		v := c.Get()
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

		time.Sleep(time.Millisecond)
		//c.notifier.Wait()
	}
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
