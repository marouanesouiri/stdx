package optional

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type state int8

const (
	stateAbsent state = iota
	statePresent
	stateNil
)

// Option is a container for an optional value of type T.
// It provides a type-safe way to represent values that may or may not be present,
// avoiding the need for nil pointers or special "zero" sentinel values.
//
// IMPORTANT:
//  - When using Option in a struct that will be marshaled to JSON,
//    always use the `json:",omitzero"` (Go 1.24+) tag.
//    This allows the Option to internally decide whether to omit the field (None), 
//    set it to "null" (Nil), or set it to the real value (Some).
type Option[T any] struct {
	state state
	value T
}

// Some creates an Option with a present value.
func Some[T any](value T) Option[T] {
	return Option[T]{
		state: statePresent,
		value: value,
	}
}

// None creates an Option with an absent value.
func None[T any]() Option[T] {
	return Option[T]{state: stateAbsent}
}

// Nil creates an Option representing an explicit "null" value.
// It behaves like None (IsPresent() == false), but when marshaled to JSON, 
// it renders as 'null' instead of being omitted.
func Nil[T any]() Option[T] {
	return Option[T]{state: stateNil}
}

// FromPtr creates an Option from a pointer.
// If the pointer is nil, it returns None, otherwise Some with the dereferenced value.
func FromPtr[T any](ptr *T) Option[T] {
	if ptr == nil {
		return None[T]()
	}
	return Some(*ptr)
}

// FromZero creates an Option from a value, treating zero values as None.
func FromZero[T comparable](value T) Option[T] {
	var zero T
	if value == zero {
		return None[T]()
	}
	return Some(value)
}

// FromPair creates an Option from a (value, ok) pair.
// If ok is false, it returns None, otherwise Some(value).
func FromPair[T any](value T, ok bool) Option[T] {
	if !ok {
		return None[T]()
	}
	return Some(value)
}

// IsPresent returns true if the value is present.
func (o Option[T]) IsPresent() bool {
	return o.state == statePresent
}

// IsAbsent returns true if the value is absent.
func (o Option[T]) IsAbsent() bool {
	return o.state != stateAbsent
}

// Get returns the value. Note that this returns the value even if absent.
// Use MustGet if you want to panic on absent values, or OrElse/OrEmpty for safe defaults.
func (o Option[T]) Get() T {
	return o.value
}

// MustGet returns the value or panics if absent.
func (o Option[T]) MustGet() T {
	if o.state != statePresent {
		panic("no element to get from option")
	}
	return o.value
}

// GetErr returns the value and an error if absent.
func (o Option[T]) GetErr() (T, error) {
	var zero T
	if o.state == stateNil {
		return zero, errors.New("value is nil")
	}
	if o.state == stateAbsent {
		return zero, errors.New("value is absent")
	}
	return o.value, nil
}

// OrElse returns the value if present, otherwise returns fallback.
func (o Option[T]) OrElse(fallback T) T {
	if o.state != statePresent {
		return fallback
	}
	return o.value
}

// OrElseGet returns the value if present, otherwise returns the value given by the supplier.
func (o Option[T]) OrElseGet(supplier func() T) T {
	if o.state != statePresent {
		return supplier()
	}
	return o.value
}

// OrElseErr returns the value if present, otherwise returns the error given by the supplier.
func (o Option[T]) OrElseErr(supplier func() error) (T, error) {
	if o.state != statePresent {
		var zero T
		return zero, supplier()
	}
	return o.value, nil
}

// OrEmpty returns the value if present, otherwise returns the zero value of T.
func (o Option[T]) OrEmpty() T {
	if o.state != statePresent {
		var empty T
		return empty
	}
	return o.value
}

// OrPanic returns the value if present, otherwise panics with the given message.
func (o Option[T]) OrPanic(message string) T {
	if o.state != statePresent {
		panic(message)
	}
	return o.value
}

// Ptr returns a pointer to the value if present, otherwise nil.
func (o Option[T]) Ptr() *T {
	if o.state != statePresent {
		return nil
	}
	return &o.value
}

// IsZero returns true if the option value is the zero value.
func (o Option[T]) IsZero() bool {
	return o.state == stateAbsent
}

// IfPresent executes the given function if the value is present.
func (o Option[T]) IfPresent(fn func(T)) {
	if o.state == statePresent {
		fn(o.value)
	}
}

// IfPresentOrElse executes the first function if the value is present,
// otherwise executes the second function.
func (o Option[T]) IfPresentOrElse(ifPresent func(T), orElse func()) {
	if o.state == statePresent {
		ifPresent(o.value)
	} else {
		orElse()
	}
}

// Filter returns the Option if it is present and the predicate returns true,
// otherwise returns None.
func (o Option[T]) Filter(predicate func(T) bool) Option[T] {
	if o.state != statePresent {
		return o
	}
	if predicate(o.value) {
		return o
	}
	return None[T]()
}

// Map transforms the value inside the Option using the provided function.
// If the Option is None, returns None.
func Map[T, U any](o Option[T], mapper func(T) U) Option[U] {
	if o.state != statePresent {
		return None[U]()
	}
	return Some(mapper(o.value))
}

// FlatMap transforms the value inside the Option using the provided function
// that returns an Option. If the Option is None, returns None.
func FlatMap[T, U any](o Option[T], mapper func(T) Option[U]) Option[U] {
	if o.state != statePresent {
		return None[U]()
	}
	return mapper(o.value)
}

// And returns None if the Option is None, otherwise returns other.
func (o Option[T]) And(other Option[T]) Option[T] {
	if o.state != statePresent {
		return None[T]()
	}
	return other
}

// Or returns the Option if it contains a value, otherwise returns other.
func (o Option[T]) Or(other Option[T]) Option[T] {
	if o.state != statePresent {
		return o
	}
	return other
}

// OrElseOption returns the Option if it contains a value,
// otherwise returns the Option provided by the supplier.
func (o Option[T]) OrElseOption(supplier func() Option[T]) Option[T] {
	if o.state != statePresent {
		return o
	}
	return supplier()
}

// Xor returns Some if exactly one of the two Options is Some, otherwise None.
func (o Option[T]) Xor(other Option[T]) Option[T] {
	if o.state == statePresent && other.state != statePresent {
		return o
	}
	if o.state != statePresent && other.state == statePresent {
		return other
	}
	return None[T]()
}

// Equal compares two Options using the provided equality function.
func (o Option[T]) Equal(other Option[T], eq func(T, T) bool) bool {
	if o.state != other.state {
		return false
	}
	return eq(o.value, other.value)
}

// MarshalJSON implements json.Marshaler.
// Present values are marshaled as their JSON representation.
// Absent values are marshaled as null.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.state == statePresent {
		return json.Marshal(o.value)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// null values are unmarshaled as None.
// All other values are unmarshaled as Some.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(bytes.ToLower(data)) == "null" {
		o.state = stateNil
		return nil
	}

	if err := json.Unmarshal(data, &o.value); err != nil {
		return err
	}

	o.state = statePresent
	return nil
}

// String returns a string representation of the Option.
func (o Option[T]) String() string {
	switch o.state {
	case stateAbsent:
		return "None"
	case stateNil:
		return "Nil"
	default:
		return fmt.Sprintf("Some(%v)", o.value)
	}
}

// TryMap applies the mapper function and returns Some(result) if no error occurs,
// otherwise returns None.
func TryMap[T, U any](o Option[T], mapper func(T) (U, error)) Option[U] {
	if o.state != statePresent {
		return None[U]()
	}
	result, err := mapper(o.value)
	if err != nil {
		return None[U]()
	}
	return Some(result)
}
