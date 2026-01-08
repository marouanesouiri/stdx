package deque

import "fmt"

const (
	// MinCapacity is the minimum capacity of the deque.
	// It is set to 16 to avoid frequent resizing for small queues.
	MinCapacity = 16
)

// Deque is a linear collection that supports element insertion and removal at
// both ends.
//
// This internal implementation uses a resizable ring buffer.
//
// This Deque implementation is not thread-safe.
type Deque[T any] struct {
	buf  []T
	head int
	tail int
	len  int
	mask int
}

// New creates a new Deque with the specified initial capacity.
// If the specified initial capacity is less than MinCapacity, the capacity is
// set to MinCapacity. The internal buffer will be allocated with a capacity
// that is a power of 2 greater than or equal to the specified capacity.
func New[T any](initialCap int) Deque[T] {
	if initialCap < MinCapacity {
		initialCap = MinCapacity
	}

	cap := 1
	for cap < initialCap {
		cap <<= 1
	}

	return Deque[T]{
		buf:  make([]T, cap),
		mask: cap - 1,
	}
}

// Len returns the number of elements in this deque.
func (d *Deque[T]) Len() int {
	return d.len
}

// Cap returns the current capacity of the deque's underlying buffer.
func (d *Deque[T]) Cap() int {
	return len(d.buf)
}

// PushBack inserts the specified element at the end of this deque.
// The capacity of the deque is automatically increased if necessary.
func (d *Deque[T]) PushBack(val T) {
	if d.len == len(d.buf) {
		d.grow()
	}

	d.buf[d.tail] = val
	d.tail = (d.tail + 1) & d.mask
	d.len++
}

// PushFront inserts the specified element at the front of this deque.
// The capacity of the deque is automatically increased if necessary.
func (d *Deque[T]) PushFront(val T) {
	if d.len == len(d.buf) {
		d.grow()
	}

	d.head = (d.head - 1) & d.mask
	d.buf[d.head] = val
	d.len++
}

// PopFront removes and returns the first element of this deque.
// Returns false if this deque is empty.
func (d *Deque[T]) PopFront() (T, bool) {
	if d.len == 0 {
		var zero T
		return zero, false
	}

	val := d.buf[d.head]

	var zero T
	d.buf[d.head] = zero

	d.head = (d.head + 1) & d.mask
	d.len--

	d.shrink()

	return val, true
}

// PopBack removes and returns the last element of this deque.
// Returns false if this deque is empty.
func (d *Deque[T]) PopBack() (T, bool) {
	if d.len == 0 {
		var zero T
		return zero, false
	}

	d.tail = (d.tail - 1) & d.mask
	val := d.buf[d.tail]

	var zero T
	d.buf[d.tail] = zero

	d.len--

	d.shrink()

	return val, true
}

// Front retrieves, but does not remove, the first element of this deque.
// Returns false if this deque is empty.
func (d *Deque[T]) Front() (T, bool) {
	if d.len == 0 {
		var zero T
		return zero, false
	}
	return d.buf[d.head], true
}

// Back retrieves, but does not remove, the last element of this deque.
// Returns false if this deque is empty.
func (d *Deque[T]) Back() (T, bool) {
	if d.len == 0 {
		var zero T
		return zero, false
	}
	idx := (d.tail - 1) & d.mask
	return d.buf[idx], true
}

// grow doubles the capacity of the deque.
func (d *Deque[T]) grow() {
	newCap := len(d.buf) << 1
	d.resize(newCap)
}

// shrink reduces the capacity of the deque if the number of elements
// falls below a certain threshold to conserve memory.
func (d *Deque[T]) shrink() {
	if len(d.buf) > MinCapacity && d.len*4 <= len(d.buf) {
		d.resize(len(d.buf) >> 1)
	}
}

// resize resizes the underlying buffer to the specified capacity.
func (d *Deque[T]) resize(newCap int) {
	newBuf := make([]T, newCap)

	if d.len > 0 {
		if d.head < d.tail {
			copy(newBuf, d.buf[d.head:d.tail])
		} else {
			n := copy(newBuf, d.buf[d.head:])
			copy(newBuf[n:], d.buf[:d.tail])
		}
	}

	d.buf = newBuf
	d.head = 0
	d.tail = d.len
	d.mask = newCap - 1
}

// Clear removes all of the elements from this deque.
// The deque will be empty after this call returns.
func (d *Deque[T]) Clear() {
	var zero T
	for i := range len(d.buf) {
		d.buf[i] = zero
	}
	d.head = 0
	d.tail = 0
	d.len = 0

	if len(d.buf) > MinCapacity {
		d.buf = make([]T, MinCapacity)
		d.mask = MinCapacity - 1
	}
}

// String returns a string representation of this deque.
func (d *Deque[T]) String() string {
	return fmt.Sprintf("Deque{len=%d, cap=%d, head=%d, tail=%d}", d.len, len(d.buf), d.head, d.tail)
}
