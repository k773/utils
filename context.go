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
