package utils

import (
	"context"
	"sync"
	"time"
)

type RateLimiter struct {
	t *time.Ticker
}

func NewRateLimiter(d time.Duration) *RateLimiter {
	return &RateLimiter{t: time.NewTicker(d)}
}

func (r *RateLimiter) Wait() {
	<-r.t.C
	return
}

func (r *RateLimiter) Done() {
	r.t.Stop()
}

type RateLimiterV2 struct {
	t         *time.Ticker
	ctx       context.Context
	wait      chan struct{}
	triggered bool
	cond      sync.Cond
}

func NewRateLimiterV2(ctx context.Context, resetTime time.Duration, limit int) *RateLimiterV2 {
	rl := &RateLimiterV2{
		t:         time.NewTicker(resetTime),
		ctx:       ctx,
		wait:      make(chan struct{}, limit),
		triggered: false,
		cond:      sync.Cond{L: &sync.Mutex{}},
	}
	go rl.Run()
	return rl
}

func (r *RateLimiterV2) Wait() {
	r.cond.L.Lock()
	if r.triggered {
		r.cond.Wait()
		r.cond.L.Unlock()
		return
	}
	r.cond.L.Unlock()
	r.wait <- struct{}{}
	return
}

// Done can be called when the rate limiter is no longer needed, but it is not mandatory. It will be called after context is done.
func (r *RateLimiterV2) Done() {
	r.t.Stop()
	close(r.wait)
	r.cond.Broadcast()
}

func (r *RateLimiterV2) Trigger() {
	r.cond.L.Lock()
	r.triggered = true
	r.cond.L.Unlock()
}

func (r *RateLimiterV2) Run() {
	for range r.t.C {
		if r.ctx.Err() != nil {
			break
		}
		for n := 0; n < len(r.wait); n++ {
			<-r.wait
		}

		r.cond.L.Lock()
		if r.triggered {
			r.triggered = false
			r.cond.Broadcast()
		}
		r.cond.L.Unlock()
	}

}

/*
	Rate limiter v3: does not require Done() to be called, but may use more cpu time.
*/

type RateLimiterV3 struct {
	l sync.Mutex

	lastRequest time.Time
	b           time.Duration // time between two Wait() releases
}

func NewRateLimiterV3(betweenTwoReleases time.Duration) *RateLimiterV3 {
	return &RateLimiterV3{
		b: betweenTwoReleases,
	}
}

func (r *RateLimiterV3) Wait() {
	r.l.Lock()
	defer r.l.Unlock()

	var t0 = time.Now()
	var td = r.b - Clamp(t0.Sub(r.lastRequest), 0, r.b)
	if !r.lastRequest.IsZero() {
		time.Sleep(td)
	}
	r.lastRequest = t0.Add(td)
}
