package utils

import "time"

type LongActionDetector struct {
	maxDuration time.Duration
	c           chan struct{}
}

func (l *LongActionDetector) Tick(timeout func()) {
	var t = time.NewTicker(l.maxDuration)
	go func() {
		select {
		case <-t.C:
			close(l.c)
			timeout()
		case <-l.c:
		}
	}()
	t.Stop()
}

func (l *LongActionDetector) Done() {
	close(l.c)
}
