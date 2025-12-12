package collectors

import (
	"strings"

	"github.com/marouanesouiri/stdx/optional"
	"github.com/marouanesouiri/stdx/set"
)

// Collector defines how to collect stream elements into a result.
// T is the element type, A is the mutable accumulator type, and R is the final result type.
type Collector[T, A, R any] interface {
	Supplier() A
	Accumulator(acc A, elem T) A
	Finisher(acc A) R
}

type sliceCollector[T any] struct{}

func (c sliceCollector[T]) Supplier() []T {
	return make([]T, 0)
}

func (c sliceCollector[T]) Accumulator(acc []T, elem T) []T {
	return append(acc, elem)
}

func (c sliceCollector[T]) Finisher(acc []T) []T {
	return acc
}

// ToSlice returns a Collector that accumulates elements into a slice.
func ToSlice[T any]() Collector[T, []T, []T] {
	return sliceCollector[T]{}
}

type setCollector[T comparable] struct{}

func (c setCollector[T]) Supplier() set.Set[T] {
	return set.New[T]()
}

func (c setCollector[T]) Accumulator(acc set.Set[T], elem T) set.Set[T] {
	acc.Add(elem)
	return acc
}

func (c setCollector[T]) Finisher(acc set.Set[T]) set.Set[T] {
	return acc
}

// ToSet returns a Collector that accumulates elements into a Set.
func ToSet[T comparable]() Collector[T, set.Set[T], set.Set[T]] {
	return setCollector[T]{}
}

type joiningCollector struct {
	separator string
	prefix    string
	suffix    string
}

func (c joiningCollector) Supplier() *strings.Builder {
	return &strings.Builder{}
}

func (c joiningCollector) Accumulator(acc *strings.Builder, elem string) *strings.Builder {
	if acc.Len() > 0 {
		acc.WriteString(c.separator)
	}
	acc.WriteString(elem)
	return acc
}

func (c joiningCollector) Finisher(acc *strings.Builder) string {
	if c.prefix == "" && c.suffix == "" {
		return acc.String()
	}
	var result strings.Builder
	result.WriteString(c.prefix)
	result.WriteString(acc.String())
	result.WriteString(c.suffix)
	return result.String()
}

// Joining returns a Collector that concatenates strings with a separator.
func Joining(separator string) Collector[string, *strings.Builder, string] {
	return joiningCollector{separator: separator}
}

// JoiningWith returns a Collector that concatenates strings with separator, prefix, and suffix.
func JoiningWith(separator, prefix, suffix string) Collector[string, *strings.Builder, string] {
	return joiningCollector{separator: separator, prefix: prefix, suffix: suffix}
}

type countingCollector[T any] struct{}

func (c countingCollector[T]) Supplier() int64 {
	return 0
}

func (c countingCollector[T]) Accumulator(acc int64, elem T) int64 {
	return acc + 1
}

func (c countingCollector[T]) Finisher(acc int64) int64 {
	return acc
}

// Counting returns a Collector that counts the number of elements.
func Counting[T any]() Collector[T, int64, int64] {
	return countingCollector[T]{}
}

type summingCollector[T any, N Number] struct {
	mapper func(T) N
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func (c summingCollector[T, N]) Supplier() N {
	return 0
}

func (c summingCollector[T, N]) Accumulator(acc N, elem T) N {
	return acc + c.mapper(elem)
}

func (c summingCollector[T, N]) Finisher(acc N) N {
	return acc
}

// Summing returns a Collector that sums numeric values extracted by the mapper.
func Summing[T any, N Number](mapper func(T) N) Collector[T, N, N] {
	return summingCollector[T, N]{mapper: mapper}
}

type avgState struct {
	sum   float64
	count int64
}

type averagingCollector[T any] struct {
	mapper func(T) float64
}

func (c averagingCollector[T]) Supplier() avgState {
	return avgState{}
}

func (c averagingCollector[T]) Accumulator(acc avgState, elem T) avgState {
	acc.sum += c.mapper(elem)
	acc.count++
	return acc
}

func (c averagingCollector[T]) Finisher(acc avgState) float64 {
	if acc.count == 0 {
		return 0
	}
	return acc.sum / float64(acc.count)
}

// Averaging returns a Collector that computes the average of numeric values.
func Averaging[T any](mapper func(T) float64) Collector[T, avgState, float64] {
	return averagingCollector[T]{mapper: mapper}
}

type minByCollector[T any] struct {
	less func(T, T) bool
}

func (c minByCollector[T]) Supplier() optional.Option[T] {
	return optional.None[T]()
}

func (c minByCollector[T]) Accumulator(acc optional.Option[T], elem T) optional.Option[T] {
	if acc.IsAbsent() {
		return optional.Some(elem)
	}
	if c.less(elem, acc.Get()) {
		return optional.Some(elem)
	}
	return acc
}

func (c minByCollector[T]) Finisher(acc optional.Option[T]) optional.Option[T] {
	return acc
}

// MinBy returns a Collector that finds the minimum element according to the less function.
func MinBy[T any](less func(T, T) bool) Collector[T, optional.Option[T], optional.Option[T]] {
	return minByCollector[T]{less: less}
}

type maxByCollector[T any] struct {
	less func(T, T) bool
}

func (c maxByCollector[T]) Supplier() optional.Option[T] {
	return optional.None[T]()
}

func (c maxByCollector[T]) Accumulator(acc optional.Option[T], elem T) optional.Option[T] {
	if acc.IsAbsent() {
		return optional.Some(elem)
	}
	if c.less(acc.Get(), elem) {
		return optional.Some(elem)
	}
	return acc
}

func (c maxByCollector[T]) Finisher(acc optional.Option[T]) optional.Option[T] {
	return acc
}

// MaxBy returns a Collector that finds the maximum element according to the less function.
func MaxBy[T any](less func(T, T) bool) Collector[T, optional.Option[T], optional.Option[T]] {
	return maxByCollector[T]{less: less}
}

type groupingByCollector[T any, K comparable] struct {
	keyFn func(T) K
}

func (c groupingByCollector[T, K]) Supplier() map[K][]T {
	return make(map[K][]T)
}

func (c groupingByCollector[T, K]) Accumulator(acc map[K][]T, elem T) map[K][]T {
	key := c.keyFn(elem)
	acc[key] = append(acc[key], elem)
	return acc
}

func (c groupingByCollector[T, K]) Finisher(acc map[K][]T) map[K][]T {
	return acc
}

// GroupingBy returns a Collector that groups elements by a key function.
func GroupingBy[T any, K comparable](keyFn func(T) K) Collector[T, map[K][]T, map[K][]T] {
	return groupingByCollector[T, K]{keyFn: keyFn}
}

type partitionState[T any] struct {
	trueList  []T
	falseList []T
}

type partitioningByCollector[T any] struct {
	predicate func(T) bool
}

func (c partitioningByCollector[T]) Supplier() partitionState[T] {
	return partitionState[T]{
		trueList:  make([]T, 0),
		falseList: make([]T, 0),
	}
}

func (c partitioningByCollector[T]) Accumulator(acc partitionState[T], elem T) partitionState[T] {
	if c.predicate(elem) {
		acc.trueList = append(acc.trueList, elem)
	} else {
		acc.falseList = append(acc.falseList, elem)
	}
	return acc
}

func (c partitioningByCollector[T]) Finisher(acc partitionState[T]) map[bool][]T {
	return map[bool][]T{
		true:  acc.trueList,
		false: acc.falseList,
	}
}

// PartitioningBy returns a Collector that partitions elements by a predicate.
func PartitioningBy[T any](predicate func(T) bool) Collector[T, partitionState[T], map[bool][]T] {
	return partitioningByCollector[T]{predicate: predicate}
}

type toMapCollector[T any, K comparable, V any] struct {
	keyFn   func(T) K
	valueFn func(T) V
	merger  func(V, V) V
}

func (c toMapCollector[T, K, V]) Supplier() map[K]V {
	return make(map[K]V)
}

func (c toMapCollector[T, K, V]) Accumulator(acc map[K]V, elem T) map[K]V {
	key := c.keyFn(elem)
	value := c.valueFn(elem)
	if c.merger != nil {
		if existing, exists := acc[key]; exists {
			acc[key] = c.merger(existing, value)
			return acc
		}
	}
	acc[key] = value
	return acc
}

func (c toMapCollector[T, K, V]) Finisher(acc map[K]V) map[K]V {
	return acc
}

// ToMap returns a Collector that collects elements into a map.
func ToMap[T any, K comparable, V any](keyFn func(T) K, valueFn func(T) V) Collector[T, map[K]V, map[K]V] {
	return toMapCollector[T, K, V]{keyFn: keyFn, valueFn: valueFn}
}

// ToMapWith returns a Collector that collects elements into a map with a merge function for duplicate keys.
func ToMapWith[T any, K comparable, V any](keyFn func(T) K, valueFn func(T) V, merger func(V, V) V) Collector[T, map[K]V, map[K]V] {
	return toMapCollector[T, K, V]{keyFn: keyFn, valueFn: valueFn, merger: merger}
}

type statsState struct {
	count int64
	sum   float64
	min   float64
	max   float64
}

// Statistics contains statistical information about a stream of values.
type Statistics struct {
	Count   int64
	Sum     float64
	Min     float64
	Max     float64
	Average float64
}

type summarizingCollector[T any] struct {
	mapper func(T) float64
}

func (c summarizingCollector[T]) Supplier() statsState {
	return statsState{}
}

func (c summarizingCollector[T]) Accumulator(acc statsState, elem T) statsState {
	value := c.mapper(elem)
	if acc.count == 0 {
		acc.min = value
		acc.max = value
	} else {
		if value < acc.min {
			acc.min = value
		}
		if value > acc.max {
			acc.max = value
		}
	}
	acc.sum += value
	acc.count++
	return acc
}

func (c summarizingCollector[T]) Finisher(acc statsState) Statistics {
	avg := 0.0
	if acc.count > 0 {
		avg = acc.sum / float64(acc.count)
	}
	return Statistics{
		Count:   acc.count,
		Sum:     acc.sum,
		Min:     acc.min,
		Max:     acc.max,
		Average: avg,
	}
}

// Summarizing returns a Collector that computes statistics for numeric values.
func Summarizing[T any](mapper func(T) float64) Collector[T, statsState, Statistics] {
	return summarizingCollector[T]{mapper: mapper}
}
