package maps

import (
	"github.com/k773/utils"
	"golang.org/x/net/context"
	"time"
)

type TimeoutMapItem[T any] struct {
	Item     T
	Deadline time.Time
}

// TimeoutMap
// You are free to use either SafeMap either M
type TimeoutMap[K comparable, T any] struct {
	*utils.SafeMap[K, *TimeoutMapItem[T]]
	OnDelete func(key K)
}

func NewTimeoutMap[K comparable, T any](onDelete func(key K)) *TimeoutMap[K, T] {
	return &TimeoutMap[K, T]{SafeMap: utils.NewSafeMap[K, *TimeoutMapItem[T]](), OnDelete: onDelete}
}

func (m *TimeoutMap[K, T]) Run(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		if ctx.Err() != nil {
			break
		}

		m.Lock()

		for k, item := range m.M {
			if t.After(item.Deadline) {
				if m.OnDelete != nil {
					m.OnDelete(k)
				}

				delete(m.M, k)
			}
		}

		m.Unlock()
	}
}
