package common

import "time"

type TimedValue[T any] struct {
	Time  time.Time `json:"time"`
	Value T         `json:"value"`
}

func (t *TimedValue[T]) Set(v T) {
	t.Time, t.Value = time.Now(), v
}
