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
	l.c = make(chan struct{})
	var t = time.NewTimer(l.maxDuration)
	go func() {
		select {
		case <-t.C:
			timeout()
		case <-l.c:
		}
		t.Stop()
	}()
}

func (l *LongActionDetector) Done() {
	close(l.c)
}
