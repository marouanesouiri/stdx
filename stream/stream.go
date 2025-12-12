package stream

import (
	"iter"
	"sort"

	"github.com/marouanesouiri/stdx/collectors"
	"github.com/marouanesouiri/stdx/optional"
)

// Stream wraps an iter.Seq and provides functional operations on sequences of elements.
// Streams support lazy evaluation and can be created from various sources including
// slices, maps, channels, and custom types implementing Streamer.
type Stream[T any] struct {
	seq iter.Seq[T]
}

// Streamer is an interface for types that can produce streams.
// Custom types can implement this interface to work seamlessly with the stream API.
type Streamer[T any] interface {
	Stream() Stream[T]
}

// From creates a Stream from a slice.
// The resulting stream will collect back to a slice by default.
func From[T any](slice []T) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			for _, v := range slice {
				if !yield(v) {
					return
				}
			}
		},
	}
}

// FromSeq creates a Stream from an iter.Seq.
// This allows integration with Go 1.23's standard iteration patterns.
func FromSeq[T any](seq iter.Seq[T]) Stream[T] {
	return Stream[T]{seq: seq}
}

// FromStreamer creates a Stream from a type implementing Streamer.
func FromStreamer[T any](s Streamer[T]) Stream[T] {
	return s.Stream()
}

// FromChannel creates a Stream from a channel.
// The stream will consume values from the channel until it is closed.
func FromChannel[T any](ch <-chan T) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			for v := range ch {
				if !yield(v) {
					return
				}
			}
		},
	}
}

// Of creates a Stream from variadic arguments.
func Of[T any](values ...T) Stream[T] {
	return From(values)
}

// Empty creates an empty Stream.
func Empty[T any]() Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {},
	}
}

// Range creates a Stream of integers from start (inclusive) to end (exclusive).
func Range(start, end int) Stream[int] {
	return Stream[int]{
		seq: func(yield func(int) bool) {
			for i := start; i < end; i++ {
				if !yield(i) {
					return
				}
			}
		},
	}
}

// Generate creates an infinite Stream by repeatedly calling the supplier function.
func Generate[T any](supplier func() T) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			for {
				if !yield(supplier()) {
					return
				}
			}
		},
	}
}

// Iterate creates an infinite Stream by applying a function to a seed value iteratively.
func Iterate[T any](seed T, fn func(T) T) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			current := seed
			for {
				if !yield(current) {
					return
				}
				current = fn(current)
			}
		},
	}
}

// Seq returns the underlying iter.Seq for use with for-range loops.
func (s Stream[T]) Seq() iter.Seq[T] {
	return s.seq
}

// Filter returns a Stream containing only elements matching the predicate.
func (s Stream[T]) Filter(predicate func(T) bool) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			for v := range s.seq {
				if predicate(v) {
					if !yield(v) {
						return
					}
				}
			}
		},
	}
}

// Map transforms each element using the mapper function.
func (s Stream[T]) Map(mapper func(T) T) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			for v := range s.seq {
				if !yield(mapper(v)) {
					return
				}
			}
		},
	}
}

// MapTo transforms each element using the mapper function and changes the element type.
func MapTo[T, U any](s Stream[T], mapper func(T) U) Stream[U] {
	return Stream[U]{
		seq: func(yield func(U) bool) {
			for v := range s.seq {
				if !yield(mapper(v)) {
					return
				}
			}
		},
	}
}

// FlatMap transforms each element to a Stream and flattens the results.
func (s Stream[T]) FlatMap(mapper func(T) Stream[T]) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			for v := range s.seq {
				for u := range mapper(v).seq {
					if !yield(u) {
						return
					}
				}
			}
		},
	}
}

// FlatMapTo transforms each element to a Stream of a different type and flattens the results.
func FlatMapTo[T, U any](s Stream[T], mapper func(T) Stream[U]) Stream[U] {
	return Stream[U]{
		seq: func(yield func(U) bool) {
			for v := range s.seq {
				for u := range mapper(v).seq {
					if !yield(u) {
						return
					}
				}
			}
		},
	}
}

// Distinct returns a Stream with duplicate elements removed.
// T must be a comparable type. For non-comparable types, use DistinctBy.
func (s Stream[T]) Distinct() Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			seen := make(map[any]struct{})
			for v := range s.seq {
				if _, exists := seen[v]; !exists {
					seen[v] = struct{}{}
					if !yield(v) {
						return
					}
				}
			}
		},
	}
}

// DistinctBy returns a Stream with duplicates removed based on a key function.
func (s Stream[T]) DistinctBy(keyFn func(T) any) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			seen := make(map[any]struct{})
			for v := range s.seq {
				key := keyFn(v)
				if _, exists := seen[key]; !exists {
					seen[key] = struct{}{}
					if !yield(v) {
						return
					}
				}
			}
		},
	}
}

// Sorted returns a Stream with elements sorted according to the less function.
// This operation materializes the entire stream into memory.
// Uses Go's standard library sort.Slice for optimal performance.
func (s Stream[T]) Sorted(less func(T, T) bool) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			slice := s.ToSlice()
			sort.Slice(slice, func(i, j int) bool {
				return less(slice[i], slice[j])
			})
			for _, v := range slice {
				if !yield(v) {
					return
				}
			}
		},
	}
}

// SortedWith returns a Stream with elements sorted using a custom sorting function.
// The sortFn should sort the provided slice in-place.
// This allows using any sorting algorithm or sort.Sort with custom types.
func (s Stream[T]) SortedWith(sortFn func([]T)) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			slice := s.ToSlice()
			sortFn(slice)
			for _, v := range slice {
				if !yield(v) {
					return
				}
			}
		},
	}
}

// Peek performs an action on each element without modifying the stream.
// Useful for debugging or side effects.
func (s Stream[T]) Peek(action func(T)) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			for v := range s.seq {
				action(v)
				if !yield(v) {
					return
				}
			}
		},
	}
}

// Limit returns a Stream with at most n elements.
func (s Stream[T]) Limit(n int64) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			count := int64(0)
			for v := range s.seq {
				if count >= n {
					return
				}
				if !yield(v) {
					return
				}
				count++
			}
		},
	}
}

// Skip returns a Stream that skips the first n elements.
func (s Stream[T]) Skip(n int64) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			count := int64(0)
			for v := range s.seq {
				if count < n {
					count++
					continue
				}
				if !yield(v) {
					return
				}
			}
		},
	}
}

// TakeWhile returns a Stream that takes elements while the predicate is true.
func (s Stream[T]) TakeWhile(predicate func(T) bool) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			for v := range s.seq {
				if !predicate(v) {
					return
				}
				if !yield(v) {
					return
				}
			}
		},
	}
}

// DropWhile returns a Stream that drops elements while the predicate is true.
func (s Stream[T]) DropWhile(predicate func(T) bool) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			dropping := true
			for v := range s.seq {
				if dropping && predicate(v) {
					continue
				}
				dropping = false
				if !yield(v) {
					return
				}
			}
		},
	}
}

// Concat returns a Stream that concatenates this stream with another.
func (s Stream[T]) Concat(other Stream[T]) Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			for v := range s.seq {
				if !yield(v) {
					return
				}
			}
			for v := range other.seq {
				if !yield(v) {
					return
				}
			}
		},
	}
}

// Reverse returns a Stream with elements in reverse order.
// This operation materializes the entire stream into memory.
func (s Stream[T]) Reverse() Stream[T] {
	return Stream[T]{
		seq: func(yield func(T) bool) {
			slice := s.ToSlice()
			for i := len(slice) - 1; i >= 0; i-- {
				if !yield(slice[i]) {
					return
				}
			}
		},
	}
}

// ForEach executes an action for each element in the stream.
func (s Stream[T]) ForEach(action func(T)) {
	for v := range s.seq {
		action(v)
	}
}

// Collect gathers stream elements using the provided Collector.
// Returns the result type R as specified by the collector.
// The return type is automatically inferred from the collector's type parameters.
func (s Stream[T]) Collect(collector collectors.Collector[T, any, any]) any {
	acc := collector.Supplier()
	for v := range s.seq {
		acc = collector.Accumulator(acc, v)
	}
	return collector.Finisher(acc)
}

// CollectTo gathers stream elements using the provided Collector with full type safety.
// Use this when you need the exact return type instead of any.
func CollectTo[T, A, R any](s Stream[T], collector collectors.Collector[T, A, R]) R {
	acc := collector.Supplier()
	for v := range s.seq {
		acc = collector.Accumulator(acc, v)
	}
	return collector.Finisher(acc)
}

// ToSlice collects all elements into a slice.
func (s Stream[T]) ToSlice() []T {
	result := make([]T, 0)
	for v := range s.seq {
		result = append(result, v)
	}
	return result
}

// Reduce combines all elements using the accumulator function, starting with identity.
func (s Stream[T]) Reduce(identity T, accumulator func(T, T) T) T {
	result := identity
	for v := range s.seq {
		result = accumulator(result, v)
	}
	return result
}

// ReduceOptional combines all elements using the accumulator function.
// Returns None if the stream is empty.
func (s Stream[T]) ReduceOptional(accumulator func(T, T) T) optional.Option[T] {
	var result T
	first := true
	for v := range s.seq {
		if first {
			result = v
			first = false
		} else {
			result = accumulator(result, v)
		}
	}
	if first {
		return optional.None[T]()
	}
	return optional.Some(result)
}

// Count returns the number of elements in the stream.
func (s Stream[T]) Count() int64 {
	count := int64(0)
	for range s.seq {
		count++
	}
	return count
}

// AnyMatch returns true if any element matches the predicate.
func (s Stream[T]) AnyMatch(predicate func(T) bool) bool {
	for v := range s.seq {
		if predicate(v) {
			return true
		}
	}
	return false
}

// AllMatch returns true if all elements match the predicate.
func (s Stream[T]) AllMatch(predicate func(T) bool) bool {
	for v := range s.seq {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// NoneMatch returns true if no elements match the predicate.
func (s Stream[T]) NoneMatch(predicate func(T) bool) bool {
	for v := range s.seq {
		if predicate(v) {
			return false
		}
	}
	return true
}

// FindFirst returns the first element wrapped in Option.
// Returns None if the stream is empty.
func (s Stream[T]) FindFirst() optional.Option[T] {
	for v := range s.seq {
		return optional.Some(v)
	}
	return optional.None[T]()
}

// FindAny returns any element from the stream.
// For sequential streams, this is equivalent to FindFirst.
func (s Stream[T]) FindAny() optional.Option[T] {
	return s.FindFirst()
}

// Min returns the minimum element according to the less function.
func (s Stream[T]) Min(less func(T, T) bool) optional.Option[T] {
	var min T
	first := true
	for v := range s.seq {
		if first || less(v, min) {
			min = v
			first = false
		}
	}
	if first {
		return optional.None[T]()
	}
	return optional.Some(min)
}

// Max returns the maximum element according to the less function.
func (s Stream[T]) Max(less func(T, T) bool) optional.Option[T] {
	var max T
	first := true
	for v := range s.seq {
		if first || less(max, v) {
			max = v
			first = false
		}
	}
	if first {
		return optional.None[T]()
	}
	return optional.Some(max)
}

// ToMap collects elements into a map using key and value functions.
func (s Stream[T]) ToMap(keyFn func(T) any, valueFn func(T) any) map[any]any {
	result := make(map[any]any)
	for v := range s.seq {
		result[keyFn(v)] = valueFn(v)
	}
	return result
}

// ToMapBy collects elements into a map using a key function, with the element as value.
func (s Stream[T]) ToMapBy(keyFn func(T) any) map[any]T {
	result := make(map[any]T)
	for v := range s.seq {
		result[keyFn(v)] = v
	}
	return result
}

// GroupBy groups elements by a key function.
func (s Stream[T]) GroupBy(keyFn func(T) any) map[any][]T {
	result := make(map[any][]T)
	for v := range s.seq {
		key := keyFn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// PartitionBy partitions elements into two slices based on a predicate.
// Returns (matching, notMatching).
func (s Stream[T]) PartitionBy(predicate func(T) bool) ([]T, []T) {
	matching := make([]T, 0)
	notMatching := make([]T, 0)
	for v := range s.seq {
		if predicate(v) {
			matching = append(matching, v)
		} else {
			notMatching = append(notMatching, v)
		}
	}
	return matching, notMatching
}
