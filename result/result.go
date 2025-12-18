package result

import (
	"fmt"

	"github.com/marouanesouiri/stdx/optional"
)

// Result represents the result of an operation that can either succeed (Ok) or fail (Err).
type Result[T any] struct {
	value T
	err   error
}

// Ok creates a successful Result containing a value.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, err: nil}
}

// Err creates a failed Result containing an error.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// From creates a Result from a standard Go (value, error) pair.
func From[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// IsOk returns true if the Result is successful.
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr returns true if the Result is a failure.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Unwrap returns the value if the Result is Ok, or panics if it is Err.
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic(fmt.Sprintf("called Result.Unwrap() on an Err value: %v", r.err))
	}
	return r.value
}

// UnwrapOr returns the value if the Result is Ok, otherwise returns the default value.
func (r Result[T]) UnwrapOr(defaultVal T) T {
	if r.err != nil {
		return defaultVal
	}
	return r.value
}

// UnwrapOrElse returns the value if the Result is Ok, otherwise computes and returns a default value.
func (r Result[T]) UnwrapOrElse(fn func() T) T {
	if r.err != nil {
		return fn()
	}
	return r.value
}

// Value returns the value stored in the Result.
// Note: It returns the zero value if the Result is Err. Check IsOk() first.
func (r Result[T]) Value() T {
	return r.value
}

// Err returns the error stored in the Result, or nil if it is Ok.
func (r Result[T]) Err() error {
	return r.err
}

// ToPair returns the (value, error) pair, making it easy to integrate with standard Go interfaces.
func (r Result[T]) ToPair() (T, error) {
	return r.value, r.err
}

// String returns a string representation of the Result.
func (r Result[T]) String() string {
	if r.err != nil {
		return fmt.Sprintf("Err(%v)", r.err)
	}
	return fmt.Sprintf("Ok(%v)", r.value)
}

// Ptr returns a pointer to the value if Ok, otherwise returns nil.
func (r Result[T]) Ptr() *T {
	if r.err != nil {
		return nil
	}
	return &r.value
}

// IfOk executes the given function if the Result is Ok.
func (r Result[T]) IfOk(fn func(T)) {
	if r.err == nil {
		fn(r.value)
	}
}

// IfErr executes the given function if the Result is Err.
func (r Result[T]) IfErr(fn func(error)) {
	if r.err != nil {
		fn(r.err)
	}
}

// IfOkOrElse executes okFn if Ok, otherwise executes errFn.
func (r Result[T]) IfOkOrElse(okFn func(T), errFn func(error)) {
	if r.err == nil {
		okFn(r.value)
	} else {
		errFn(r.err)
	}
}

// Or returns the Result if Ok, otherwise returns the other Result.
func (r Result[T]) Or(other Result[T]) Result[T] {
	if r.err == nil {
		return r
	}
	return other
}

// OrElseResult computes and returns a fallback Result if the original is Err.
func (r Result[T]) OrElseResult(fn func() Result[T]) Result[T] {
	if r.err == nil {
		return r
	}
	return fn()
}

// Filter checks the value with a predicate. If the Result is Ok and the predicate returns false,
// it returns an Err with the provided error. Otherwise it returns the original Result.
func (r Result[T]) Filter(predicate func(T) bool, ifFalseErr error) Result[T] {
	if r.err != nil {
		return r
	}
	if !predicate(r.value) {
		return Err[T](ifFalseErr)
	}
	return r
}

// Option converts the Result to an optional.Option.
// If Result is Ok, it returns Some(value). If Err, it returns None.
func (r Result[T]) Option() optional.Option[T] {
	if r.err != nil {
		return optional.None[T]()
	}
	return optional.Some(r.value)
}

// Recover returns the value if Ok, otherwise handles the error with the provided function and returns its result.
func (r Result[T]) Recover(fn func(error) T) T {
	if r.err != nil {
		return fn(r.err)
	}
	return r.value
}

// Void is a Result that contains no value.
// It is used for operations that can fail but don't return data on success.
type Void struct {
	err error
}

// OkVoid creates a successful Result with no value.
func OkVoid() Void {
	return Void{err: nil}
}

// ErrVoid creates a failed Result with no value.
func ErrVoid(err error) Void {
	return Void{err: err}
}

// IsOk returns true if the Result is successful.
func (v Void) IsOk() bool {
	return v.err == nil
}

// IsErr returns true if the Result is a failure.
func (v Void) IsErr() bool {
	return v.err != nil
}

// Err creates a failed Result containing an error.
func (v Void) Err() error {
	return v.err
}

// ToResult converts a Void into a Result[struct{}].
func (v Void) ToResult() Result[struct{}] {
	if v.err != nil {
		return Err[struct{}](v.err)
	}
	return Ok(struct{}{})
}

// String returns a string representation of the Result.
func (v Void) String() string {
	if v.err != nil {
		return fmt.Sprintf("Err(%v)", v.err)
	}
	return fmt.Sprintf("Ok")
}
