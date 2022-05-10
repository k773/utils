package files

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

type RateLimitedWriter struct {
	writer   io.Writer
	writerAt io.WriterAt
	closer   io.Closer

	s, fs sync.Mutex

	maxBytesPerSecond int

	ctx                context.Context
	cancel             context.CancelFunc
	canWriteThisSecond int
	lastSecondStart    time.Time
	cond               *sync.Cond
}

func NewRateLimitedWriter(f interface{}, maxBytesPerSecond int) *RateLimitedWriter {
	var r = &RateLimitedWriter{
		maxBytesPerSecond:  maxBytesPerSecond,
		canWriteThisSecond: maxBytesPerSecond,
		lastSecondStart:    time.Now(),
		cond:               sync.NewCond(&sync.Mutex{}),
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.writer, _ = f.(io.Writer)
	r.writerAt, _ = f.(io.WriterAt)
	r.closer, _ = f.(io.Closer)
	go r.condUpdater()
	return r
}

func (r *RateLimitedWriter) condUpdater() {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for range t.C {
		if r.ctx.Err() != nil {
			break
		}
		r.s.Lock()
		r.canWriteThisSecond = r.maxBytesPerSecond
		r.s.Unlock()
		r.cond.Broadcast()
	}
}

func (r *RateLimitedWriter) Close() {
	r.cancel()
	if r.closer != nil {
		_ = r.closer.Close()
	}
}

func (r *RateLimitedWriter) Write(data []byte) (wrote int, e error) {
	if r.writer == nil {
		return 0, errors.New("interface does not implement io.Closer")
	}

	var toWrite int
	var n int

	r.fs.Lock()
	r.cond.L.Lock()
	for e == nil && wrote < len(data) {
		r.s.Lock()
		toWrite = len(data) - wrote
		if toWrite > r.canWriteThisSecond {
			toWrite = r.canWriteThisSecond
		}
		r.s.Unlock()
		if toWrite != 0 {
			fmt.Printf("wrote: %v, toWrite: %v, sum:%v\n", wrote, toWrite, toWrite+wrote)
			n, e = r.writer.Write(data[wrote : wrote+toWrite])
		}
		if e == nil {
			r.s.Lock()
			wrote += n
			if wrote < len(data) {
				r.s.Unlock()
				r.cond.Wait()
			} else {
				r.s.Unlock()
			}
		}
	}
	r.cond.L.Unlock()
	r.fs.Unlock()
	return
}

func (r *RateLimitedWriter) WriteAt(data []byte, offset int64) (wrote int, e error) {
	if r.writerAt == nil {
		return 0, errors.New("interface does not implement io.Closer")
	}

	var toWrite int
	var n int

	//r.fs.Lock()
	r.cond.L.Lock()
	for e == nil && wrote < len(data) {
		r.s.Lock()
		toWrite = len(data) - wrote
		if toWrite > r.canWriteThisSecond {
			toWrite = r.canWriteThisSecond
		}
		r.canWriteThisSecond -= toWrite
		r.s.Unlock()
		if toWrite != 0 {
			//fmt.Printf("wrote: %v, toWrite: %v, sum:%v\n", wrote, toWrite, toWrite+wrote)
			n, e = r.writerAt.WriteAt(data[wrote:wrote+toWrite], offset+int64(wrote))
		}
		if e == nil {
			r.s.Lock()
			wrote += n
			if wrote < len(data) {
				r.s.Unlock()
				r.cond.Wait()
			} else {
				r.s.Unlock()
			}
		}
	}
	r.cond.L.Unlock()
	//r.fs.Unlock()
	return
}
