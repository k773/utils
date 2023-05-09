package utils

import (
	"golang.org/x/net/context"
	"sync"
	"time"
)

func RunForeverFirstSyncThenAsync(ctx context.Context, f func(), every time.Duration, simultaneously bool) {
	var run = func() {
		if simultaneously {
			go f()
		} else {
			f()
		}
	}

	if ctx.Err() == nil {
		run()
	} else {
		return
	}

	go func() {
		var t = time.NewTicker(every)
		defer t.Stop()
	l:
		for {
			select {
			case <-ctx.Done():
				break l
			case <-t.C:
				run()
			}
		}
	}()
}

func RunForeverAsyncNoFirstTime(ctx context.Context, f func(), every time.Duration, simultaneously bool) {
	var run = func() {
		if simultaneously {
			go f()
		} else {
			f()
		}
	}

	if ctx.Err() == nil {
		go func() {
			var t = time.NewTicker(every)
			defer t.Stop()
		l:
			for {
				select {
				case <-ctx.Done():
					break l
				case <-t.C:
					run()
				}
			}
		}()
	}
}

func RunForeverSync(ctx context.Context, f func(), every time.Duration, simultaneously bool) {
	var run = func() {
		if simultaneously {
			go f()
		} else {
			f()
		}
	}

	if ctx.Err() == nil {
		run()
	} else {
		return
	}

	var t = time.NewTicker(every)
	defer t.Stop()
l:
	for {
		select {
		case <-ctx.Done():
			break l
		case <-t.C:
			run()
		}
	}
}

// RunForeverSyncUntil : f func should return false to stop the execution
func RunForeverSyncUntil(ctx context.Context, f func() bool, every time.Duration) {
	if ctx.Err() != nil || !f() {
		return
	}

	var t = time.NewTicker(every)
	defer t.Stop()
l:
	for {
		select {
		case <-ctx.Done():
			break l
		case <-t.C:
			if !f() {
				break l
			}
		}
	}
}

// DelayedExecution will trigger wg only at the end of the execution; increasing the value before calling this function is up to the caller
func DelayedExecution(ctx context.Context, wg *sync.WaitGroup, executeOnParentCancelled bool, delay time.Duration, f func()) {
	var wait = time.NewTimer(delay)
	defer wait.Stop()
	if wg != nil {
		defer wg.Done()
	}

	select {
	case <-wait.C:
		f()
	case <-ctx.Done():
		if executeOnParentCancelled {
			f()
		}
	}
}

func ExecuteWitWGAsync(w *sync.WaitGroup, f func()) {
	w.Add(1)
	go func() {
		defer w.Done()
		f()
	}()
}

func WithLock(s sync.Locker, f func()) func() {
	return func() {
		s.Lock()
		defer s.Unlock()
	}
}
