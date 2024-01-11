package synctools

type ChanMutex struct {
	c chan struct{}
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
