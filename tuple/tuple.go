package tuple

import (
	"encoding/json"
	"fmt"
)

// Tuple2 represents a pair of two values of potentially different types.
// It provides a type-safe way to group two related values together.
type Tuple2[T, U any] struct {
	First  T
	Second U
}

// NewTuple2 creates a new Tuple2 with the given values.
func NewTuple2[T, U any](first T, second U) Tuple2[T, U] {
	return Tuple2[T, U]{
		First:  first,
		Second: second,
	}
}

// FromPair is an alias for NewTuple2 for backward compatibility.
func FromPair[T, U any](first T, second U) Tuple2[T, U] {
	return NewTuple2(first, second)
}

// Values returns both values as individual return values.
func (t Tuple2[T, U]) Values() (T, U) {
	return t.First, t.Second
}

// Swap returns a new Tuple2 with the first and second values swapped.
func (t Tuple2[T, U]) Swap() Tuple2[U, T] {
	return Tuple2[U, T]{
		First:  t.Second,
		Second: t.First,
	}
}

// MapFirst applies the given function to the first value and returns a new Tuple2.
func (t Tuple2[T, U]) MapFirst(fn func(T) T) Tuple2[T, U] {
	return Tuple2[T, U]{
		First:  fn(t.First),
		Second: t.Second,
	}
}

// MapSecond applies the given function to the second value and returns a new Tuple2.
func (t Tuple2[T, U]) MapSecond(fn func(U) U) Tuple2[T, U] {
	return Tuple2[T, U]{
		First:  t.First,
		Second: fn(t.Second),
	}
}

// Map applies the given function to both values and returns a new Tuple2.
func (t Tuple2[T, U]) Map(fn func(T, U) (T, U)) Tuple2[T, U] {
	first, second := fn(t.First, t.Second)
	return Tuple2[T, U]{
		First:  first,
		Second: second,
	}
}

// MapBoth applies different functions to each value and returns a new Tuple2.
func MapBoth[T, U, V, W any](t Tuple2[T, U], fnFirst func(T) V, fnSecond func(U) W) Tuple2[V, W] {
	return Tuple2[V, W]{
		First:  fnFirst(t.First),
		Second: fnSecond(t.Second),
	}
}

// Equal compares two Tuple2 values using the provided equality functions.
func (t Tuple2[T, U]) Equal(other Tuple2[T, U], eqT func(T, T) bool, eqU func(U, U) bool) bool {
	return eqT(t.First, other.First) && eqU(t.Second, other.Second)
}

// String returns a string representation of the Tuple2.
func (t Tuple2[T, U]) String() string {
	return fmt.Sprintf("(%v, %v)", t.First, t.Second)
}

// MarshalJSON implements json.Marshaler.
// The tuple is marshaled as a JSON array with two elements.
func (t Tuple2[T, U]) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{t.First, t.Second})
}

// UnmarshalJSON implements json.Unmarshaler.
// Expects a JSON array with exactly two elements.
func (t *Tuple2[T, U]) UnmarshalJSON(data []byte) error {
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) != 2 {
		return fmt.Errorf("expected array of length 2, got %d", len(arr))
	}
	if err := json.Unmarshal(arr[0], &t.First); err != nil {
		return err
	}
	if err := json.Unmarshal(arr[1], &t.Second); err != nil {
		return err
	}
	return nil
}

// Tuple3 represents a triple of three values of potentially different types.
type Tuple3[T, U, V any] struct {
	First  T
	Second U
	Third  V
}

// NewTuple3 creates a new Tuple3 with the given values.
func NewTuple3[T, U, V any](first T, second U, third V) Tuple3[T, U, V] {
	return Tuple3[T, U, V]{
		First:  first,
		Second: second,
		Third:  third,
	}
}

// Values returns all three values as individual return values.
func (t Tuple3[T, U, V]) Values() (T, U, V) {
	return t.First, t.Second, t.Third
}

// MapFirst applies the given function to the first value and returns a new Tuple3.
func (t Tuple3[T, U, V]) MapFirst(fn func(T) T) Tuple3[T, U, V] {
	return Tuple3[T, U, V]{
		First:  fn(t.First),
		Second: t.Second,
		Third:  t.Third,
	}
}

// MapSecond applies the given function to the second value and returns a new Tuple3.
func (t Tuple3[T, U, V]) MapSecond(fn func(U) U) Tuple3[T, U, V] {
	return Tuple3[T, U, V]{
		First:  t.First,
		Second: fn(t.Second),
		Third:  t.Third,
	}
}

// MapThird applies the given function to the third value and returns a new Tuple3.
func (t Tuple3[T, U, V]) MapThird(fn func(V) V) Tuple3[T, U, V] {
	return Tuple3[T, U, V]{
		First:  t.First,
		Second: t.Second,
		Third:  fn(t.Third),
	}
}

// Map applies the given function to all three values and returns a new Tuple3.
func (t Tuple3[T, U, V]) Map(fn func(T, U, V) (T, U, V)) Tuple3[T, U, V] {
	first, second, third := fn(t.First, t.Second, t.Third)
	return Tuple3[T, U, V]{
		First:  first,
		Second: second,
		Third:  third,
	}
}

// String returns a string representation of the Tuple3.
func (t Tuple3[T, U, V]) String() string {
	return fmt.Sprintf("(%v, %v, %v)", t.First, t.Second, t.Third)
}

// MarshalJSON implements json.Marshaler.
// The tuple is marshaled as a JSON array with three elements.
func (t Tuple3[T, U, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{t.First, t.Second, t.Third})
}

// UnmarshalJSON implements json.Unmarshaler.
// Expects a JSON array with exactly three elements.
func (t *Tuple3[T, U, V]) UnmarshalJSON(data []byte) error {
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) != 3 {
		return fmt.Errorf("expected array of length 3, got %d", len(arr))
	}
	if err := json.Unmarshal(arr[0], &t.First); err != nil {
		return err
	}
	if err := json.Unmarshal(arr[1], &t.Second); err != nil {
		return err
	}
	if err := json.Unmarshal(arr[2], &t.Third); err != nil {
		return err
	}
	return nil
}

// Tuple4 represents a quadruple of four values of potentially different types.
type Tuple4[T, U, V, W any] struct {
	First  T
	Second U
	Third  V
	Fourth W
}

// NewTuple4 creates a new Tuple4 with the given values.
func NewTuple4[T, U, V, W any](first T, second U, third V, fourth W) Tuple4[T, U, V, W] {
	return Tuple4[T, U, V, W]{
		First:  first,
		Second: second,
		Third:  third,
		Fourth: fourth,
	}
}

// Values returns all four values as individual return values.
func (t Tuple4[T, U, V, W]) Values() (T, U, V, W) {
	return t.First, t.Second, t.Third, t.Fourth
}

// String returns a string representation of the Tuple4.
func (t Tuple4[T, U, V, W]) String() string {
	return fmt.Sprintf("(%v, %v, %v, %v)", t.First, t.Second, t.Third, t.Fourth)
}

// MarshalJSON implements json.Marshaler.
// The tuple is marshaled as a JSON array with four elements.
func (t Tuple4[T, U, V, W]) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{t.First, t.Second, t.Third, t.Fourth})
}

// UnmarshalJSON implements json.Unmarshaler.
// Expects a JSON array with exactly four elements.
func (t *Tuple4[T, U, V, W]) UnmarshalJSON(data []byte) error {
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) != 4 {
		return fmt.Errorf("expected array of length 4, got %d", len(arr))
	}
	if err := json.Unmarshal(arr[0], &t.First); err != nil {
		return err
	}
	if err := json.Unmarshal(arr[1], &t.Second); err != nil {
		return err
	}
	if err := json.Unmarshal(arr[2], &t.Third); err != nil {
		return err
	}
	if err := json.Unmarshal(arr[3], &t.Fourth); err != nil {
		return err
	}
	return nil
}

// Tuple5 represents a quintuple of five values of potentially different types.
type Tuple5[T, U, V, W, X any] struct {
	First  T
	Second U
	Third  V
	Fourth W
	Fifth  X
}

// NewTuple5 creates a new Tuple5 with the given values.
func NewTuple5[T, U, V, W, X any](first T, second U, third V, fourth W, fifth X) Tuple5[T, U, V, W, X] {
	return Tuple5[T, U, V, W, X]{
		First:  first,
		Second: second,
		Third:  third,
		Fourth: fourth,
		Fifth:  fifth,
	}
}

// Values returns all five values as individual return values.
func (t Tuple5[T, U, V, W, X]) Values() (T, U, V, W, X) {
	return t.First, t.Second, t.Third, t.Fourth, t.Fifth
}

// String returns a string representation of the Tuple5.
func (t Tuple5[T, U, V, W, X]) String() string {
	return fmt.Sprintf("(%v, %v, %v, %v, %v)", t.First, t.Second, t.Third, t.Fourth, t.Fifth)
}

// MarshalJSON implements json.Marshaler.
// The tuple is marshaled as a JSON array with five elements.
func (t Tuple5[T, U, V, W, X]) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{t.First, t.Second, t.Third, t.Fourth, t.Fifth})
}

// UnmarshalJSON implements json.Unmarshaler.
// Expects a JSON array with exactly five elements.
func (t *Tuple5[T, U, V, W, X]) UnmarshalJSON(data []byte) error {
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) != 5 {
		return fmt.Errorf("expected array of length 5, got %d", len(arr))
	}
	if err := json.Unmarshal(arr[0], &t.First); err != nil {
		return err
	}
	if err := json.Unmarshal(arr[1], &t.Second); err != nil {
		return err
	}
	if err := json.Unmarshal(arr[2], &t.Third); err != nil {
		return err
	}
	if err := json.Unmarshal(arr[3], &t.Fourth); err != nil {
		return err
	}
	if err := json.Unmarshal(arr[4], &t.Fifth); err != nil {
		return err
	}
	return nil
}

// Zip combines two slices into a slice of Tuple2.
// The resulting slice has the length of the shorter input slice.
func Zip[T, U any](first []T, second []U) []Tuple2[T, U] {
	minLen := len(first)
	if len(second) < minLen {
		minLen = len(second)
	}
	result := make([]Tuple2[T, U], minLen)
	for i := 0; i < minLen; i++ {
		result[i] = NewTuple2(first[i], second[i])
	}
	return result
}

// Unzip splits a slice of Tuple2 into two separate slices.
func Unzip[T, U any](tuples []Tuple2[T, U]) ([]T, []U) {
	first := make([]T, len(tuples))
	second := make([]U, len(tuples))
	for i, t := range tuples {
		first[i] = t.First
		second[i] = t.Second
	}
	return first, second
}

// Zip3 combines three slices into a slice of Tuple3.
// The resulting slice has the length of the shortest input slice.
func Zip3[T, U, V any](first []T, second []U, third []V) []Tuple3[T, U, V] {
	minLen := len(first)
	if len(second) < minLen {
		minLen = len(second)
	}
	if len(third) < minLen {
		minLen = len(third)
	}
	result := make([]Tuple3[T, U, V], minLen)
	for i := 0; i < minLen; i++ {
		result[i] = NewTuple3(first[i], second[i], third[i])
	}
	return result
}

// Unzip3 splits a slice of Tuple3 into three separate slices.
func Unzip3[T, U, V any](tuples []Tuple3[T, U, V]) ([]T, []U, []V) {
	first := make([]T, len(tuples))
	second := make([]U, len(tuples))
	third := make([]V, len(tuples))
	for i, t := range tuples {
		first[i] = t.First
		second[i] = t.Second
		third[i] = t.Third
	}
	return first, second, third
}
