package optional

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
)

// Option is a container for an optional value of type T.
// It provides a type-safe way to represent values that may or may not be present,
// avoiding the need for nil pointers or special "zero" sentinel values.
type Option[T any] struct {
	isPresent bool
	value     T
}

// Some creates an Option with a present value.
func Some[T any](value T) Option[T] {
	return Option[T]{
		isPresent: true,
		value:     value,
	}
}

// None creates an Option with an absent value.
func None[T any]() Option[T] {
	return Option[T]{
		isPresent: false,
	}
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

// IsPresent returns true if the value is present.
func (o Option[T]) IsPresent() bool {
	return o.isPresent
}

// IsAbsent returns true if the value is absent.
func (o Option[T]) IsAbsent() bool {
	return !o.isPresent
}

// Get returns the value. Note that this returns the value even if absent.
// Use MustGet if you want to panic on absent values, or OrElse/OrEmpty for safe defaults.
func (o Option[T]) Get() T {
	return o.value
}

// MustGet returns the value or panics if absent.
func (o Option[T]) MustGet() T {
	if !o.isPresent {
		panic("no element to get from option")
	}
	return o.value
}

// GetErr returns the value and an error if absent.
func (o Option[T]) GetErr() (T, error) {
	if !o.isPresent {
		var zero T
		return zero, errors.New("value is absent")
	}
	return o.value, nil
}

// OrElse returns the value if present, otherwise returns fallback.
func (o Option[T]) OrElse(fallback T) T {
	if !o.isPresent {
		return fallback
	}
	return o.value
}

// OrElseGet returns the value if present, otherwise returns the value given by the supplier.
func (o Option[T]) OrElseGet(supplier func() T) T {
	if !o.isPresent {
		return supplier()
	}
	return o.value
}

// OrElseErr returns the value if present, otherwise returns the error given by the supplier.
func (o Option[T]) OrElseErr(supplier func() error) (T, error) {
	if !o.isPresent {
		var zero T
		return zero, supplier()
	}
	return o.value, nil
}

// OrEmpty returns the value if present, otherwise returns the zero value of T.
func (o Option[T]) OrEmpty() T {
	var empty T
	if !o.isPresent {
		return empty
	}
	return o.value
}

// OrPanic returns the value if present, otherwise panics with the given message.
func (o Option[T]) OrPanic(message string) T {
	if !o.isPresent {
		panic(message)
	}
	return o.value
}

// Ptr returns a pointer to the value if present, otherwise nil.
func (o Option[T]) Ptr() *T {
	if !o.isPresent {
		return nil
	}
	return &o.value
}

// IsZero returns true if the option value is the zero value.
func (o Option[T]) IsZero() bool {
	if !o.isPresent {
		return true
	}
	var v any = o.value
	if v, ok := v.(interface{ IsZero() bool }); ok {
		return v.IsZero()
	}
	return reflect.ValueOf(o.value).IsZero()
}

// IfPresent executes the given function if the value is present.
func (o Option[T]) IfPresent(fn func(T)) {
	if o.isPresent {
		fn(o.value)
	}
}

// IfPresentOrElse executes the first function if the value is present,
// otherwise executes the second function.
func (o Option[T]) IfPresentOrElse(ifPresent func(T), orElse func()) {
	if o.isPresent {
		ifPresent(o.value)
	} else {
		orElse()
	}
}

// Filter returns the Option if it is present and the predicate returns true,
// otherwise returns None.
func (o Option[T]) Filter(predicate func(T) bool) Option[T] {
	if !o.isPresent {
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
	if !o.isPresent {
		return None[U]()
	}
	return Some(mapper(o.value))
}

// FlatMap transforms the value inside the Option using the provided function
// that returns an Option. If the Option is None, returns None.
func FlatMap[T, U any](o Option[T], mapper func(T) Option[U]) Option[U] {
	if !o.isPresent {
		return None[U]()
	}
	return mapper(o.value)
}

// And returns None if the Option is None, otherwise returns other.
func (o Option[T]) And(other Option[T]) Option[T] {
	if !o.isPresent {
		return None[T]()
	}
	return other
}

// Or returns the Option if it contains a value, otherwise returns other.
func (o Option[T]) Or(other Option[T]) Option[T] {
	if o.isPresent {
		return o
	}
	return other
}

// OrElseOption returns the Option if it contains a value,
// otherwise returns the Option provided by the supplier.
func (o Option[T]) OrElseOption(supplier func() Option[T]) Option[T] {
	if o.isPresent {
		return o
	}
	return supplier()
}

// Xor returns Some if exactly one of the two Options is Some, otherwise None.
func (o Option[T]) Xor(other Option[T]) Option[T] {
	if o.isPresent && !other.isPresent {
		return o
	}
	if !o.isPresent && other.isPresent {
		return other
	}
	return None[T]()
}

// Equal compares two Options using the provided equality function.
func (o Option[T]) Equal(other Option[T], eq func(T, T) bool) bool {
	if o.isPresent != other.isPresent {
		return false
	}
	if !o.isPresent {
		return true
	}
	return eq(o.value, other.value)
}

// MarshalJSON implements json.Marshaler.
// Present values are marshaled as their JSON representation.
// Absent values are marshaled as null.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.isPresent {
		return json.Marshal(o.value)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements json.Unmarshaler.
// null values are unmarshaled as None.
// All other values are unmarshaled as Some.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(bytes.ToLower(data)) == "null" {
		o.isPresent = false
		return nil
	}

	if err := json.Unmarshal(data, &o.value); err != nil {
		return err
	}

	o.isPresent = true
	return nil
}

// String returns a string representation of the Option.
func (o Option[T]) String() string {
	if !o.isPresent {
		return "None"
	}
	return "Some(" + any(o.value).(interface{ String() string }).String() + ")"
}

// TryMap applies the mapper function and returns Some(result) if no error occurs,
// otherwise returns None.
func TryMap[T, U any](o Option[T], mapper func(T) (U, error)) Option[U] {
	if !o.isPresent {
		return None[U]()
	}
	result, err := mapper(o.value)
	if err != nil {
		return None[U]()
	}
	return Some(result)
}
