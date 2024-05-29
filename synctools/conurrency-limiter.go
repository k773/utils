package synctools

import "context"

// ConcurrencyLimiter is a tool that helps to limit the number of concurrent operations.
type ConcurrencyLimiter struct {
	limit chan struct{}
}

// NewConcurrencyLimiter creates a new ConcurrencyLimiter with the specified limit.
// The limit must be in the range of [1, +inf).
func NewConcurrencyLimiter(limit int) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{limit: make(chan struct{}, max(1, limit))}
}

func (c *ConcurrencyLimiter) Lock() {
	c.limit <- struct{}{}
}

// TryLock tries to acquire a lock with the specified context.
func (c *ConcurrencyLimiter) TryLock(ctx context.Context) bool {
	select {
	case c.limit <- struct{}{}:
		return true
	case <-ctx.Done():
		return false
	}
}

// Unlock releases the lock.
func (c *ConcurrencyLimiter) Unlock() {
	<-c.limit
}

// WithLock wraps the specified function with a lock.
func (c *ConcurrencyLimiter) WithLock(ctx context.Context, do func() error) error {
	ok := c.TryLock(ctx)
	if !ok {
		return ctx.Err()
	}
	defer c.Unlock()
	return do()
}
