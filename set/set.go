package set

import "iter"

// Set is a collection of unique comparable elements.
// It is safe to copy Set values as the underlying map is a reference type.
type Set[T comparable] struct {
	items map[T]struct{}
}

// New creates and returns a new empty Set.
func New[T comparable]() Set[T] {
	return Set[T]{
		items: make(map[T]struct{}),
	}
}

// FromSlice creates a new Set containing all unique elements from the given slice.
func FromSlice[T comparable](slice []T) Set[T] {
	s := New[T]()
	for _, item := range slice {
		s.Add(item)
	}
	return s
}

// Add inserts an element into the set.
// Returns true if the element was added (wasn't already present), false otherwise.
func (s *Set[T]) Add(item T) bool {
	if _, exists := s.items[item]; exists {
		return false
	}
	s.items[item] = struct{}{}
	return true
}

// AddAll inserts multiple elements into the set.
// Returns the count of elements that were actually added (excludes duplicates).
func (s *Set[T]) AddAll(items ...T) int {
	count := 0
	for _, item := range items {
		if s.Add(item) {
			count++
		}
	}
	return count
}

// Remove deletes an element from the set.
// Returns true if the element was removed (was present), false otherwise.
func (s *Set[T]) Remove(item T) bool {
	if _, exists := s.items[item]; !exists {
		return false
	}
	delete(s.items, item)
	return true
}

// Contains checks if an element exists in the set.
func (s *Set[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

// Size returns the number of elements in the set.
func (s *Set[T]) Size() int {
	return len(s.items)
}

// AsMap returns the underlying map representation of the Set.
// The returned map has keys of type T and values of empty struct{}.
// Useful for interoperability with code that expects map[T]struct{}.
//
// WARNING: This returns a reference to the internal map, NOT a copy.
// Modifying the returned map will affect the Set. Use Clone() if you need an independent copy.
func (s Set[T]) AsMap() map[T]struct{} {
	return s.items
}

// IsEmpty returns true if the set contains no elements.
func (s *Set[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Clear removes all elements from the set.
func (s *Set[T]) Clear() {
	for item := range s.items {
		delete(s.items, item)
	}
}

// ToSlice returns a slice containing all elements in the set.
// The order of elements is not guaranteed.
func (s *Set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s.items))
	for item := range s.items {
		slice = append(slice, item)
	}
	return slice
}

// Range calls the given function for each element in the set.
// If the function returns false, iteration stops.
func (s *Set[T]) Range(fn func(T) bool) {
	for item := range s.items {
		if !fn(item) {
			return
		}
	}
}

// Union returns a new set containing all elements from both sets.
func (s *Set[T]) Union(other Set[T]) Set[T] {
	result := New[T]()
	for item := range s.items {
		result.Add(item)
	}
	for item := range other.items {
		result.Add(item)
	}
	return result
}

// Intersection returns a new set containing only elements present in both sets.
func (s *Set[T]) Intersection(other Set[T]) Set[T] {
	result := New[T]()
	for item := range s.items {
		if other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// Difference returns a new set containing elements in this set but not in the other set.
func (s *Set[T]) Difference(other Set[T]) Set[T] {
	result := New[T]()
	for item := range s.items {
		if !other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// SymmetricDifference returns a new set containing elements in either set but not in both.
func (s *Set[T]) SymmetricDifference(other Set[T]) Set[T] {
	result := New[T]()
	for item := range s.items {
		if !other.Contains(item) {
			result.Add(item)
		}
	}
	for item := range other.items {
		if !s.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// IsSubset returns true if all elements of this set are in the other set.
func (s *Set[T]) IsSubset(other Set[T]) bool {
	for item := range s.items {
		if !other.Contains(item) {
			return false
		}
	}
	return true
}

// IsSuperset returns true if this set contains all elements of the other set.
func (s *Set[T]) IsSuperset(other Set[T]) bool {
	return other.IsSubset(*s)
}

// Equal returns true if both sets contain exactly the same elements.
func (s *Set[T]) Equal(other Set[T]) bool {
	if s.Size() != other.Size() {
		return false
	}
	for item := range s.items {
		if !other.Contains(item) {
			return false
		}
	}
	return true
}

// Clone creates a deep copy of the Set with an independent internal map.
// Modifications to the clone will not affect the original Set and vice versa.
func (s *Set[T]) Clone() Set[T] {
	clone := New[T]()
	for item := range s.items {
		clone.Add(item)
	}
	return clone
}

// Seq returns an iter.Seq that yields all elements in the set.
// This enables use with Go 1.23 for-range loops.
func (s Set[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for item := range s.items {
			if !yield(item) {
				return
			}
		}
	}
}
