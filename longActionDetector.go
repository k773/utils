package utils

import "time"

type LongActionDetector struct {
	maxDuration time.Duration
	c           chan struct{}
}

func NewLongActionDetector(timeout time.Duration) *LongActionDetector {
	return &LongActionDetector{maxDuration: timeout}
}

func (l *LongActionDetector) Tick(timeout func()) {
	var t = time.NewTimer(l.maxDuration)
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
