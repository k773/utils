package utils

import (
	"context"
	"time"
)

func WithContextTimeout(parent context.Context, timeout time.Duration, f func(ctx context.Context)) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	f(ctx)
}

func WithContextDeadline(parent context.Context, deadline time.Time, f func(ctx context.Context)) {
	ctx, cancel := context.WithDeadline(parent, deadline)
	defer cancel()
	f(ctx)
}

//
//type MergedContext struct {
//	Ctx1 context.Context
//	Ctx2 context.Context
//}
//
//func (m *MergedContext) Deadline() (deadline time.Time, ok bool) {
//	d1, ok1 := m.Ctx1.Deadline()
//	d2, ok2 := m.Ctx2.Deadline()
//
//	if ok1 && ok2 {
//		if d1.Before(d2) {
//			return d1, true
//		} else {
//			return d2, true
//		}
//	} else if ok1 {
//		return d1, true
//	} else {
//		return time.Time{}, false
//	}
//}
//
//func (m *MergedContext) Done() <-chan struct{} {
//
//}
//
//func (m *MergedContext) Err() error        {}
//func (m *MergedContext) Value(key any) any {}
