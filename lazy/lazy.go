package lazy

import (
	"sync"
)

// Lazy represents a value that is computed only once, on first access.
// It provides thread-safe lazy initialization using sync.Once.
// Multiple goroutines can safely call Get() - the supplier function will only execute once.
type Lazy[T any] struct {
	once     sync.Once
	supplier func() T
	value    T
	computed bool
}

// New creates a new Lazy value that will compute its value using the supplier function
// when first accessed via Get(). The supplier will only be called once.
func New[T any](supplier func() T) Lazy[T] {
	return Lazy[T]{
		supplier: supplier,
	}
}

// Of creates a Lazy value that is already computed with the given value.
// No supplier function is needed, and Get() will immediately return the value.
func Of[T any](value T) Lazy[T] {
	return Lazy[T]{
		value:    value,
		computed: true,
	}
}

// Get forces the computation if not already done and returns the value.
// This method is thread-safe - if multiple goroutines call Get() concurrently,
// the supplier function will only execute once.
func (l *Lazy[T]) Get() T {
	l.once.Do(func() {
		if l.supplier != nil {
			l.value = l.supplier()
		}
		l.computed = true
	})
	return l.value
}

// IsComputed returns true if the value has been computed, false otherwise.
// This method is safe to call concurrently with Get().
func (l *Lazy[T]) IsComputed() bool {
	isComputed := false
	l.once.Do(func() {
		if l.supplier != nil {
			l.value = l.supplier()
		}
		l.computed = true
		isComputed = true
	})
	return !isComputed || l.computed
}

// Map creates a new Lazy value by applying the transformation function to this Lazy value.
// The transformation is also lazy - it won't execute until the returned Lazy is accessed.
func Map[T, U any](l *Lazy[T], fn func(T) U) Lazy[U] {
	return New(func() U {
		return fn(l.Get())
	})
}

// FlatMap creates a new Lazy value by applying a function that returns a Lazy value.
// This flattens the nested Lazy values. Both computations are lazy.
func FlatMap[T, U any](l *Lazy[T], fn func(T) *Lazy[U]) Lazy[U] {
	return New(func() U {
		innerLazy := fn(l.Get())
		return innerLazy.Get()
	})
}

// Filter creates a Lazy that returns the value if the predicate is true,
// otherwise returns the zero value of T.
func Filter[T any](l *Lazy[T], predicate func(T) bool) Lazy[T] {
	return New(func() T {
		val := l.Get()
		if predicate(val) {
			return val
		}
		var zero T
		return zero
	})
}

// OrElse returns this Lazy if its value satisfies the predicate,
// otherwise returns the alternative Lazy.
func OrElse[T any](l *Lazy[T], predicate func(T) bool, alternative *Lazy[T]) Lazy[T] {
	return New(func() T {
		val := l.Get()
		if predicate(val) {
			return val
		}
		return alternative.Get()
	})
}
