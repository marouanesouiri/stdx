package stream

import (
	"strconv"
	"testing"
)

func TestFrom(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	s := From(slice)
	result := s.ToSlice()
	if len(result) != 5 {
		t.Errorf("expected 5 elements, got %d", len(result))
	}
}

func TestFilter(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5, 6})
	result := s.Filter(func(x int) bool { return x%2 == 0 }).ToSlice()
	if len(result) != 3 {
		t.Errorf("expected 3 elements, got %d", len(result))
	}
	for _, v := range result {
		if v%2 != 0 {
			t.Errorf("expected even number, got %d", v)
		}
	}
}

func TestMap(t *testing.T) {
	s := From([]int{1, 2, 3})
	result := s.Map(func(x int) int { return x * 2 }).ToSlice()
	expected := []int{2, 4, 6}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("expected %d, got %d", expected[i], v)
		}
	}
}

func TestFlatMap(t *testing.T) {
	s := From([]int{1, 2, 3})
	result := s.FlatMap(func(x int) Stream[int] {
		return Of(x, x*10)
	}).ToSlice()
	expected := []int{1, 10, 2, 20, 3, 30}
	if len(result) != len(expected) {
		t.Errorf("expected %d elements, got %d", len(expected), len(result))
	}
}

func TestDistinct(t *testing.T) {
	s := From([]int{1, 2, 2, 3, 3, 3, 4})
	result := s.Distinct().ToSlice()
	if len(result) != 4 {
		t.Errorf("expected 4 unique elements, got %d", len(result))
	}
}

func TestDistinctBy(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 30},
	}
	s := From(people)
	result := s.DistinctBy(func(p Person) any { return p.Age }).ToSlice()
	if len(result) != 2 {
		t.Errorf("expected 2 unique ages, got %d", len(result))
	}
}

func TestSorted(t *testing.T) {
	s := From([]int{5, 2, 8, 1, 9})
	result := s.Sorted(func(a, b int) bool { return a < b }).ToSlice()
	for i := 1; i < len(result); i++ {
		if result[i] < result[i-1] {
			t.Errorf("not sorted: %v", result)
		}
	}
}

func TestLimit(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5})
	result := s.Limit(3).ToSlice()
	if len(result) != 3 {
		t.Errorf("expected 3 elements, got %d", len(result))
	}
}

func TestSkip(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5})
	result := s.Skip(2).ToSlice()
	if len(result) != 3 {
		t.Errorf("expected 3 elements, got %d", len(result))
	}
	if result[0] != 3 {
		t.Errorf("expected first element to be 3, got %d", result[0])
	}
}

func TestTakeWhile(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5})
	result := s.TakeWhile(func(x int) bool { return x < 4 }).ToSlice()
	if len(result) != 3 {
		t.Errorf("expected 3 elements, got %d", len(result))
	}
}

func TestDropWhile(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5})
	result := s.DropWhile(func(x int) bool { return x < 3 }).ToSlice()
	if len(result) != 3 {
		t.Errorf("expected 3 elements, got %d", len(result))
	}
	if result[0] != 3 {
		t.Errorf("expected first element to be 3, got %d", result[0])
	}
}

func TestConcat(t *testing.T) {
	s1 := From([]int{1, 2, 3})
	s2 := From([]int{4, 5, 6})
	result := s1.Concat(s2).ToSlice()
	if len(result) != 6 {
		t.Errorf("expected 6 elements, got %d", len(result))
	}
}

func TestReverse(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5})
	result := s.Reverse().ToSlice()
	expected := []int{5, 4, 3, 2, 1}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("expected %d, got %d", expected[i], v)
		}
	}
}

func TestReduce(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5})
	result := s.Reduce(0, func(a, b int) int { return a + b })
	if result != 15 {
		t.Errorf("expected 15, got %d", result)
	}
}

func TestReduceOptional(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5})
	result := s.ReduceOptional(func(a, b int) int { return a + b })
	if result.IsAbsent() || result.Get() != 15 {
		t.Errorf("expected Some(15), got %v", result)
	}

	empty := Empty[int]()
	emptyResult := empty.ReduceOptional(func(a, b int) int { return a + b })
	if emptyResult.IsPresent() {
		t.Errorf("expected None, got %v", emptyResult)
	}
}

func TestCount(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5})
	count := s.Count()
	if count != 5 {
		t.Errorf("expected 5, got %d", count)
	}
}

func TestAnyMatch(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5})
	if !s.AnyMatch(func(x int) bool { return x > 3 }) {
		t.Error("expected true")
	}
	if s.AnyMatch(func(x int) bool { return x > 10 }) {
		t.Error("expected false")
	}
}

func TestAllMatch(t *testing.T) {
	s := From([]int{2, 4, 6, 8})
	if !s.AllMatch(func(x int) bool { return x%2 == 0 }) {
		t.Error("expected true")
	}
	s2 := From([]int{1, 2, 3})
	if s2.AllMatch(func(x int) bool { return x%2 == 0 }) {
		t.Error("expected false")
	}
}

func TestNoneMatch(t *testing.T) {
	s := From([]int{1, 3, 5, 7})
	if !s.NoneMatch(func(x int) bool { return x%2 == 0 }) {
		t.Error("expected true")
	}
}

func TestFindFirst(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5})
	result := s.FindFirst()
	if result.IsAbsent() || result.Get() != 1 {
		t.Errorf("expected Some(1), got %v", result)
	}

	empty := Empty[int]()
	emptyResult := empty.FindFirst()
	if emptyResult.IsPresent() {
		t.Errorf("expected None, got %v", emptyResult)
	}
}

func TestMinMax(t *testing.T) {
	s := From([]int{5, 2, 8, 1, 9})
	min := s.Min(func(a, b int) bool { return a < b })
	if min.IsAbsent() || min.Get() != 1 {
		t.Errorf("expected 1, got %v", min)
	}

	s2 := From([]int{5, 2, 8, 1, 9})
	max := s2.Max(func(a, b int) bool { return a < b })
	if max.IsAbsent() || max.Get() != 9 {
		t.Errorf("expected 9, got %v", max)
	}
}

func TestGroupBy(t *testing.T) {
	words := []string{"apple", "apricot", "banana", "berry", "cherry"}
	s := From(words)
	result := s.GroupBy(func(s string) any { return rune(s[0]) })
	if len(result[rune('a')]) != 2 {
		t.Errorf("expected 2 words starting with 'a', got %d", len(result[rune('a')]))
	}
	if len(result[rune('b')]) != 2 {
		t.Errorf("expected 2 words starting with 'b', got %d", len(result[rune('b')]))
	}
}

func TestPartitionBy(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5, 6})
	evens, odds := s.PartitionBy(func(x int) bool { return x%2 == 0 })
	if len(evens) != 3 {
		t.Errorf("expected 3 evens, got %d", len(evens))
	}
	if len(odds) != 3 {
		t.Errorf("expected 3 odds, got %d", len(odds))
	}
}

func TestRange(t *testing.T) {
	s := Range(1, 6)
	result := s.ToSlice()
	if len(result) != 5 {
		t.Errorf("expected 5 elements, got %d", len(result))
	}
}

func TestOf(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)
	result := s.ToSlice()
	if len(result) != 5 {
		t.Errorf("expected 5 elements, got %d", len(result))
	}
}

func TestEmpty(t *testing.T) {
	s := Empty[int]()
	result := s.ToSlice()
	if len(result) != 0 {
		t.Errorf("expected 0 elements, got %d", len(result))
	}
}

func TestLazyEvaluation(t *testing.T) {
	count := 0
	s := From([]int{1, 2, 3, 4, 5})
	s.Peek(func(x int) { count++ }).
		Filter(func(x int) bool { return x%2 == 0 })

	if count != 0 {
		t.Errorf("expected lazy evaluation, but count is %d", count)
	}

	result := s.Peek(func(x int) { count++ }).
		Filter(func(x int) bool { return x%2 == 0 }).
		ToSlice()

	if count != len([]int{1, 2, 3, 4, 5}) {
		t.Errorf("expected count to be 5 after collection, got %d", count)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 even numbers, got %d", len(result))
	}
}

func TestSeqIntegration(t *testing.T) {
	s := From([]int{1, 2, 3, 4, 5})
	sum := 0
	for v := range s.Filter(func(x int) bool { return x%2 == 0 }).Seq() {
		sum += v
	}
	if sum != 6 {
		t.Errorf("expected sum of evens to be 6, got %d", sum)
	}
}

func BenchmarkFilter(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).Filter(func(x int) bool { return x%2 == 0 }).ToSlice()
	}
}

func BenchmarkMap(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).Map(func(x int) int { return x * 2 }).ToSlice()
	}
}

func BenchmarkChainedOperations(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapTo(
			From(data).
				Filter(func(x int) bool { return x%2 == 0 }).
				Limit(100),
			func(x int) int { return x * 2 },
		).ToSlice()
	}
}

func BenchmarkDistinct(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i % 100
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).Distinct().ToSlice()
	}
}

func BenchmarkSorted(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = 1000 - i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).Sorted(func(a, b int) bool { return a < b }).ToSlice()
	}
}

func BenchmarkGroupBy(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).GroupBy(func(x int) any { return x % 10 })
	}
}

func BenchmarkReduce(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).Reduce(0, func(a, b int) int { return a + b })
	}
}

func BenchmarkComplexPipeline(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).
			Filter(func(x int) bool { return x%2 == 0 }).
			Limit(1000).
			Peek(func(x int) {}).
			ToSlice()
	}
}

func BenchmarkMapToString(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapTo(From(data), strconv.Itoa).ToSlice()
	}
}
