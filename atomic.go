package utils

import "sync/atomic"

type AtomicInt int64

func (i *AtomicInt) Get() int64 {
	return atomic.LoadInt64((*int64)(i))
}

func (i *AtomicInt) Set(v int64) {
	atomic.StoreInt64((*int64)(i), v)
}

func (i *AtomicInt) Add(v int64) {
	atomic.AddInt64((*int64)(i), v)
}
