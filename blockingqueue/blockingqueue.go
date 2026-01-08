package blockingqueue

import (
	"context"
)

// BlockingQueue is a thread-safe queue backed by a Go channel.
//
// It provides a familiar object-oriented API around standard Go channels.
// Note that this does not support unbounded capacity or PushFront,
// as channels do not support these operations.
type BlockingQueue[T any] struct {
	ch chan T
}

// New creates a new BlockingQueue with the specified capacity.
// If capacity is 0, it creates an unbuffered (synchronous) queue.
func New[T any](capacity int) *BlockingQueue[T] {
	if capacity < 0 {
		capacity = 0
	}
	return &BlockingQueue[T]{
		ch: make(chan T, capacity),
	}
}

// Push inserts the specified element into this queue, waiting if necessary
// for space to become available.
func (bq *BlockingQueue[T]) Push(val T) {
	bq.ch <- val
}

// PushCtx inserts the specified element into this queue, waiting if necessary
// for space to become available or until the context is done.
// Returns nil on success, or ctx.Err() if the context is cancelled.
func (bq *BlockingQueue[T]) PushCtx(ctx context.Context, val T) error {
	select {
	case bq.ch <- val:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Pop retrieves and removes the head of this queue, waiting if necessary
// until an element becomes available.
func (bq *BlockingQueue[T]) Pop() T {
	return <-bq.ch
}

// PopCtx retrieves and removes the head of this queue, waiting if necessary
// until an element becomes available or the context is done.
// Returns (value, nil) on success, or (zero, ctx.Err()) if the context is cancelled.
func (bq *BlockingQueue[T]) PopCtx(ctx context.Context) (T, error) {
	select {
	case val := <-bq.ch:
		return val, nil
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

// TryPush inserts the specified element into this queue if it is possible to do
// so immediately without violating capacity restrictions.
// Returns true upon success and false if no space is currently available.
func (bq *BlockingQueue[T]) TryPush(val T) bool {
	select {
	case bq.ch <- val:
		return true
	default:
		return false
	}
}

// TryPop retrieves and removes the head of this queue only if it is available.
// Returns the element and true if the queue was not empty.
func (bq *BlockingQueue[T]) TryPop() (T, bool) {
	select {
	case val := <-bq.ch:
		return val, true
	default:
		var zero T
		return zero, false
	}
}

// Len returns the number of elements in the buffer.
// Note: for an unbuffered queue, this always returns 0.
func (bq *BlockingQueue[T]) Len() int {
	return len(bq.ch)
}

// Cap returns the capacity of the buffer.
func (bq *BlockingQueue[T]) Cap() int {
	return cap(bq.ch)
}

// Clear removes all of the elements from this queue.
func (bq *BlockingQueue[T]) Clear() {
	// Drain the channel
	for {
		select {
		case <-bq.ch:
		default:
			return
		}
	}
}
