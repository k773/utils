package files

import (
	"io/ioutil"
	"sync"
	"time"
)

type Cache struct {
	s sync.Mutex
	l time.Duration
	m map[string][]byte
	t map[time.Time][]string
}

func NewFilesCache(lifetime time.Duration) *Cache {
	var c = &Cache{
		m: map[string][]byte{},
		t: map[time.Time][]string{},
		l: lifetime,
	}
	go c.autoClean()
	return c
}

func (c *Cache) Read(path string) (b []byte, e error) {
	var t = time.Now().Truncate(c.l)
	c.s.Lock()
	var h bool
	if b, h = c.m[path]; h {
		return
	}
	c.s.Unlock()

	if b, e = ioutil.ReadFile(path); e == nil {
		c.s.Lock()
		c.m[path] = b
		c.t[t] = append(c.t[t], path)
		c.s.Unlock()
	}
	return
}

func (c *Cache) autoClean() {
	for t := range time.NewTicker(c.l / 3).C {
		c.s.Lock()
		k_ := t.Truncate(c.l)
		for _, k := range c.t[k_] {
			delete(c.m, k)
		}
		delete(c.t, k_)
		c.s.Unlock()
	}
}
