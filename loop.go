package utils

import (
	"golang.org/x/net/context"
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
		for range t.C {
			if ctx.Err() != nil {
				break
			}
			run()
		}
		t.Stop()
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
			for range t.C {
				if ctx.Err() != nil {
					break
				}
				run()
			}
			t.Stop()
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
	for range t.C {
		if ctx.Err() != nil {
			break
		}
		run()
	}
	t.Stop()
}
