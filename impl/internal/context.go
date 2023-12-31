package internal

import (
	"context"
	"time"
)

// CancelFunc is a function api that can be used to stop a context.
type CancelFunc context.CancelFunc

// WithCancel returns a copy of ctx with a new Done channel.
// See context.WithCancel for more details.
func WithCancel(ctx context.Context) (context.Context, CancelFunc) {
	c, cancel := context.WithCancel(ctx)
	return c, CancelFunc(cancel)
}

// WithDeadline returns a copy of ctx with the deadline adjusted to be no later than deadline.
// See context.WithDeadline for more details.
func WithDeadline(ctx context.Context, deadline time.Time) (context.Context, CancelFunc) {
	c, cancel := context.WithDeadline(ctx, deadline)
	return c, CancelFunc(cancel)
}

// WithTimeout is shorthand for ctx.WithDeadline(time.Now().Add(duration)).
// See context.Context.WithTimeout for more details.
func WithTimeout(ctx context.Context, duration time.Duration) (context.Context, CancelFunc) {
	return WithDeadline(ctx, time.Now().Add(duration))
}

// ShouldStop returns a chan that's closed when work done on behalf of this
// context should be stopped.
// See context.Context.Done for more details.
func ShouldStop(ctx context.Context) <-chan struct{} {
	return ctx.Done()
}

// StopReason returns a non-nil error value after Done is closed.
// See context.Context.Err for more details.
func StopReason(ctx context.Context) error {
	return ctx.Err()
}

// Stopped is shorthand for StopReason(ctx) != nil because it increases the readability of common use cases.
func Stopped(ctx context.Context) bool {
	return ctx.Err() != nil
}
