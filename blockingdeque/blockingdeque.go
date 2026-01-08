package blockingdeque

import (
	"context"
	"sync"

	"github.com/marouanesouiri/stdx/deque"
)

// BlockingDeque is a thread-safe double-ended queue.
type BlockingDeque[T any] struct {
	mu       sync.Mutex
	q        deque.Deque[T]
	capacity int

	notEmpty chan struct{}
	notFull  chan struct{}
}

// New creates a new BlockingDeque with the specified capacity.
func New[T any](capacity int) *BlockingDeque[T] {
	if capacity < 1 {
		capacity = 1
	}
	bd := &BlockingDeque[T]{
		q:        deque.New[T](capacity),
		capacity: capacity,
		notEmpty: make(chan struct{}, 1),
		notFull:  make(chan struct{}, 1),
	}

	bd.notFull <- struct{}{}

	return bd
}

// PushBack inserts the specified element at the back.
func (bd *BlockingDeque[T]) PushBack(val T) {
	bd.PushBackCtx(context.Background(), val)
}

func (bd *BlockingDeque[T]) PushBackCtx(ctx context.Context, val T) error {
	for {
		bd.mu.Lock()
		if bd.q.Len() < bd.capacity {
			bd.q.PushBack(val)

			select {
			case bd.notEmpty <- struct{}{}:
			default:
			}

			if bd.q.Len() < bd.capacity {
				select {
				case bd.notFull <- struct{}{}:
				default:
				}
			}

			bd.mu.Unlock()
			return nil
		}
		bd.mu.Unlock()

		select {
		case <-bd.notFull:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// PushFront inserts the specified element at the front.
func (bd *BlockingDeque[T]) PushFront(val T) {
	bd.PushFrontCtx(context.Background(), val)
}

func (bd *BlockingDeque[T]) PushFrontCtx(ctx context.Context, val T) error {
	for {
		bd.mu.Lock()
		if bd.q.Len() < bd.capacity {
			bd.q.PushFront(val)

			select {
			case bd.notEmpty <- struct{}{}:
			default:
			}

			if bd.q.Len() < bd.capacity {
				select {
				case bd.notFull <- struct{}{}:
				default:
				}
			}
			bd.mu.Unlock()
			return nil
		}
		bd.mu.Unlock()

		select {
		case <-bd.notFull:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// PopFront retrieves and removes the first element of this deque.
func (bd *BlockingDeque[T]) PopFront() T {
	val, _ := bd.PopFrontCtx(context.Background())
	return val
}

func (bd *BlockingDeque[T]) PopFrontCtx(ctx context.Context) (T, error) {
	for {
		bd.mu.Lock()
		if bd.q.Len() > 0 {
			val, _ := bd.q.PopFront()

			select {
			case bd.notFull <- struct{}{}:
			default:
			}

			if bd.q.Len() > 0 {
				select {
				case bd.notEmpty <- struct{}{}:
				default:
				}
			}

			bd.mu.Unlock()
			return val, nil
		}
		bd.mu.Unlock()

		select {
		case <-bd.notEmpty:
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		}
	}
}

// PopBack retrieves and removes the last element of this deque.
func (bd *BlockingDeque[T]) PopBack() T {
	val, _ := bd.PopBackCtx(context.Background())
	return val
}

func (bd *BlockingDeque[T]) PopBackCtx(ctx context.Context) (T, error) {
	for {
		bd.mu.Lock()
		if bd.q.Len() > 0 {
			val, _ := bd.q.PopBack()

			select {
			case bd.notFull <- struct{}{}:
			default:
			}

			if bd.q.Len() > 0 {
				select {
				case bd.notEmpty <- struct{}{}:
				default:
				}
			}
			bd.mu.Unlock()
			return val, nil
		}
		bd.mu.Unlock()

		select {
		case <-bd.notEmpty:
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		}
	}
}

// TryPushBack inserts at the back if possible immediately.
func (bd *BlockingDeque[T]) TryPushBack(val T) bool {
	bd.mu.Lock()
	if bd.q.Len() >= bd.capacity {
		bd.mu.Unlock()
		return false
	}

	bd.q.PushBack(val)

	select {
	case bd.notEmpty <- struct{}{}:
	default:
	}

	if bd.q.Len() < bd.capacity {
		select {
		case bd.notFull <- struct{}{}:
		default:
		}
	}
	bd.mu.Unlock()
	return true
}

// TryPushFront inserts at the front if possible immediately.
func (bd *BlockingDeque[T]) TryPushFront(val T) bool {
	bd.mu.Lock()
	if bd.q.Len() >= bd.capacity {
		bd.mu.Unlock()
		return false
	}

	bd.q.PushFront(val)

	select {
	case bd.notEmpty <- struct{}{}:
	default:
	}

	if bd.q.Len() < bd.capacity {
		select {
		case bd.notFull <- struct{}{}:
		default:
		}
	}
	bd.mu.Unlock()
	return true
}

// TryPopFront retrieves from the front if available.
func (bd *BlockingDeque[T]) TryPopFront() (T, bool) {
	bd.mu.Lock()
	if bd.q.Len() == 0 {
		bd.mu.Unlock()
		var zero T
		return zero, false
	}

	val, _ := bd.q.PopFront()

	select {
	case bd.notFull <- struct{}{}:
	default:
	}

	if bd.q.Len() > 0 {
		select {
		case bd.notEmpty <- struct{}{}:
		default:
		}
	}
	bd.mu.Unlock()
	return val, true
}

// TryPopBack retrieves from the back if available.
func (bd *BlockingDeque[T]) TryPopBack() (T, bool) {
	bd.mu.Lock()
	if bd.q.Len() == 0 {
		bd.mu.Unlock()
		var zero T
		return zero, false
	}

	val, _ := bd.q.PopBack()

	select {
	case bd.notFull <- struct{}{}:
	default:
	}

	if bd.q.Len() > 0 {
		select {
		case bd.notEmpty <- struct{}{}:
		default:
		}
	}
	bd.mu.Unlock()
	return val, true
}

// Len returns the number of elements in the deque.
func (bd *BlockingDeque[T]) Len() int {
	bd.mu.Lock()
	l := bd.q.Len()
	bd.mu.Unlock()
	return l
}

// Cap returns the capacity.
func (bd *BlockingDeque[T]) Cap() int {
	return bd.capacity
}

// Clear removes all elements.
func (bd *BlockingDeque[T]) Clear() {
	bd.mu.Lock()

	bd.q.Clear()

	select {
	case <-bd.notEmpty:
	default:
	}

	select {
	case bd.notFull <- struct{}{}:
	default:
	}

	bd.mu.Unlock()
}
