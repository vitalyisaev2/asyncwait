package asyncwait

import (
	"context"
	"time"
)

// AsyncWait async wait representation
type AsyncWait interface {
	// Wait method for wait result
	Wait(func() bool) bool
}

var _ AsyncWait = (*asyncWait)(nil)

type asyncWait struct {
	pollInterval time.Duration
	timeout      time.Duration
	doneCh       chan struct{}
}

// NewAsyncWait constructor for AsyncWait
func NewAsyncWait(timeout, pollInterval time.Duration) AsyncWait {
	return &asyncWait{
		pollInterval: pollInterval,
		timeout:      timeout,
		doneCh:       make(chan struct{}),
	}
}

// Wait while timeout, make polls every pollInterval for the predicate while is not truth
func (aw asyncWait) Wait(predicate func() bool) bool {
	ctx, cancel := context.WithTimeout(context.Background(), aw.timeout)
	defer cancel()

	for {
		select {
		case <-aw.doneCh:
			return true
		case <-ctx.Done():
			return false
		case <-time.After(aw.pollInterval):
			go func() {
				if predicate() {
					select {
					case aw.doneCh <- struct{}{}:
					case <-ctx.Done():
						return
					}
				}
			}()
		}
	}
}
