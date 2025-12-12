package either

import (
	"encoding/json"
	"fmt"
)

// Either represents a value of one of two possible types (a disjoint union).
// An Either is either Left or Right.
// By convention, Left is used for failure/error cases and Right is used for success cases,
// but both are valid values and neither represents an absence of value.
type Either[L, R any] struct {
	isLeft bool
	left   L
	right  R
}

// Left creates an Either with a left value.
// By convention, left values represent failure or error cases.
func Left[L, R any](value L) Either[L, R] {
	return Either[L, R]{
		isLeft: true,
		left:   value,
	}
}

// Right creates an Either with a right value.
// By convention, right values represent success cases.
func Right[L, R any](value R) Either[L, R] {
	return Either[L, R]{
		isLeft: false,
		right:  value,
	}
}

// IsLeft returns true if this Either contains a left value.
func (e Either[L, R]) IsLeft() bool {
	return e.isLeft
}

// IsRight returns true if this Either contains a right value.
func (e Either[L, R]) IsRight() bool {
	return !e.isLeft
}

// Left returns the left value.
// Note: This returns the value even if it's a Right. Use IsLeft() to check first,
// or use LeftOr() for a safe default.
func (e Either[L, R]) Left() L {
	return e.left
}

// Right returns the right value.
// Note: This returns the value even if it's a Left. Use IsRight() to check first,
// or use RightOr() for a safe default.
func (e Either[L, R]) Right() R {
	return e.right
}

// LeftOr returns the left value if present, otherwise returns the fallback.
func (e Either[L, R]) LeftOr(fallback L) L {
	if e.isLeft {
		return e.left
	}
	return fallback
}

// RightOr returns the right value if present, otherwise returns the fallback.
func (e Either[L, R]) RightOr(fallback R) R {
	if !e.isLeft {
		return e.right
	}
	return fallback
}

// GetLeft returns the left value and a boolean indicating if it was left.
func (e Either[L, R]) GetLeft() (L, bool) {
	return e.left, e.isLeft
}

// GetRight returns the right value and a boolean indicating if it was right.
func (e Either[L, R]) GetRight() (R, bool) {
	return e.right, !e.isLeft
}

// Swap returns a new Either with left and right swapped.
func (e Either[L, R]) Swap() Either[R, L] {
	if e.isLeft {
		return Right[R, L](e.left)
	}
	return Left[R, L](e.right)
}

// MapLeft applies the function to the left value if present, otherwise returns the Either unchanged.
func (e Either[L, R]) MapLeft(fn func(L) L) Either[L, R] {
	if !e.isLeft {
		return e
	}
	return Left[L, R](fn(e.left))
}

// MapRight applies the function to the right value if present, otherwise returns the Either unchanged.
func (e Either[L, R]) MapRight(fn func(R) R) Either[L, R] {
	if e.isLeft {
		return e
	}
	return Right[L, R](fn(e.right))
}

// Map is an alias for MapRight, following the convention that Right is the "success" path.
func (e Either[L, R]) Map(fn func(R) R) Either[L, R] {
	return e.MapRight(fn)
}

// MapBothLeft applies different functions to transform the left type.
func MapBothLeft[L, R, L2 any](e Either[L, R], fn func(L) L2) Either[L2, R] {
	if e.isLeft {
		return Left[L2, R](fn(e.left))
	}
	return Right[L2, R](e.right)
}

// MapBothRight applies different functions to transform the right type.
func MapBothRight[L, R, R2 any](e Either[L, R], fn func(R) R2) Either[L, R2] {
	if !e.isLeft {
		return Right[L, R2](fn(e.right))
	}
	return Left[L, R2](e.left)
}

// MapBoth transforms both the left and right types using the provided functions.
func MapBoth[L, R, L2, R2 any](e Either[L, R], leftFn func(L) L2, rightFn func(R) R2) Either[L2, R2] {
	if e.isLeft {
		return Left[L2, R2](leftFn(e.left))
	}
	return Right[L2, R2](rightFn(e.right))
}

// Fold applies one of two functions depending on whether this is Left or Right,
// and returns the result. This is useful for pattern matching.
func Fold[L, R, T any](e Either[L, R], leftFn func(L) T, rightFn func(R) T) T {
	if e.isLeft {
		return leftFn(e.left)
	}
	return rightFn(e.right)
}

// FlatMap applies a function that returns an Either to the right value,
// flattening the result. Returns the original Either if it's a Left.
func FlatMap[L, R, R2 any](e Either[L, R], fn func(R) Either[L, R2]) Either[L, R2] {
	if e.isLeft {
		return Left[L, R2](e.left)
	}
	return fn(e.right)
}

// OrElse returns this Either if it's a Right, otherwise returns the alternative.
func (e Either[L, R]) OrElse(alternative Either[L, R]) Either[L, R] {
	if !e.isLeft {
		return e
	}
	return alternative
}

// String returns a string representation of the Either.
func (e Either[L, R]) String() string {
	if e.isLeft {
		return fmt.Sprintf("Left(%v)", e.left)
	}
	return fmt.Sprintf("Right(%v)", e.right)
}

// eitherJSON is used for JSON marshaling to include type information.
type eitherJSON struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}

// MarshalJSON implements json.Marshaler.
// The Either is marshaled as an object with "type" (either "left" or "right") and "value" fields.
func (e Either[L, R]) MarshalJSON() ([]byte, error) {
	if e.isLeft {
		return json.Marshal(eitherJSON{
			Type:  "left",
			Value: e.left,
		})
	}
	return json.Marshal(eitherJSON{
		Type:  "right",
		Value: e.right,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
// Expects a JSON object with "type" and "value" fields.
func (e *Either[L, R]) UnmarshalJSON(data []byte) error {
	var ej eitherJSON
	if err := json.Unmarshal(data, &ej); err != nil {
		return err
	}

	switch ej.Type {
	case "left":
		valueBytes, err := json.Marshal(ej.Value)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(valueBytes, &e.left); err != nil {
			return err
		}
		e.isLeft = true
	case "right":
		valueBytes, err := json.Marshal(ej.Value)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(valueBytes, &e.right); err != nil {
			return err
		}
		e.isLeft = false
	default:
		return fmt.Errorf("invalid either type: %s (expected 'left' or 'right')", ej.Type)
	}

	return nil
}
