package collectors

import (
	"testing"
)

func TestToSlice(t *testing.T) {
	collector := ToSlice[int]()
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, 1)
	acc = collector.Accumulator(acc, 2)
	acc = collector.Accumulator(acc, 3)
	acc = collector.Accumulator(acc, 4)
	acc = collector.Accumulator(acc, 5)
	result := collector.Finisher(acc)
	if len(result) != 5 {
		t.Errorf("expected 5 elements, got %d", len(result))
	}
}

func TestToSet(t *testing.T) {
	collector := ToSet[int]()
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, 1)
	acc = collector.Accumulator(acc, 2)
	acc = collector.Accumulator(acc, 2)
	acc = collector.Accumulator(acc, 3)
	acc = collector.Accumulator(acc, 3)
	acc = collector.Accumulator(acc, 3)
	result := collector.Finisher(acc)
	if result.Size() != 3 {
		t.Errorf("expected 3 unique elements, got %d", result.Size())
	}
}

func TestJoining(t *testing.T) {
	collector := Joining(", ")
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, "apple")
	acc = collector.Accumulator(acc, "banana")
	acc = collector.Accumulator(acc, "cherry")
	result := collector.Finisher(acc)
	expected := "apple, banana, cherry"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestJoiningWith(t *testing.T) {
	collector := JoiningWith(", ", "[", "]")
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, "apple")
	acc = collector.Accumulator(acc, "banana")
	acc = collector.Accumulator(acc, "cherry")
	result := collector.Finisher(acc)
	expected := "[apple, banana, cherry]"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestCounting(t *testing.T) {
	collector := Counting[int]()
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, 1)
	acc = collector.Accumulator(acc, 2)
	acc = collector.Accumulator(acc, 3)
	acc = collector.Accumulator(acc, 4)
	acc = collector.Accumulator(acc, 5)
	result := collector.Finisher(acc)
	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}

func TestSumming(t *testing.T) {
	collector := Summing(func(x int) int { return x })
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, 1)
	acc = collector.Accumulator(acc, 2)
	acc = collector.Accumulator(acc, 3)
	acc = collector.Accumulator(acc, 4)
	acc = collector.Accumulator(acc, 5)
	result := collector.Finisher(acc)
	if result != 15 {
		t.Errorf("expected 15, got %d", result)
	}
}

func TestSummingFloat(t *testing.T) {
	collector := Summing(func(x float64) float64 { return x })
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, 1.5)
	acc = collector.Accumulator(acc, 2.5)
	acc = collector.Accumulator(acc, 3.0)
	result := collector.Finisher(acc)
	expected := 7.0
	if result != expected {
		t.Errorf("expected %f, got %f", expected, result)
	}
}

func TestAveraging(t *testing.T) {
	collector := Averaging(func(x int) float64 { return float64(x) })
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, 1)
	acc = collector.Accumulator(acc, 2)
	acc = collector.Accumulator(acc, 3)
	acc = collector.Accumulator(acc, 4)
	acc = collector.Accumulator(acc, 5)
	result := collector.Finisher(acc)
	expected := 3.0
	if result != expected {
		t.Errorf("expected %f, got %f", expected, result)
	}
}

func TestMinBy(t *testing.T) {
	collector := MinBy(func(a, b int) bool { return a < b })
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, 5)
	acc = collector.Accumulator(acc, 2)
	acc = collector.Accumulator(acc, 8)
	acc = collector.Accumulator(acc, 1)
	acc = collector.Accumulator(acc, 9)
	result := collector.Finisher(acc)
	if result.IsAbsent() || result.Get() != 1 {
		t.Errorf("expected 1, got %v", result)
	}
}

func TestMaxBy(t *testing.T) {
	collector := MaxBy(func(a, b int) bool { return a < b })
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, 5)
	acc = collector.Accumulator(acc, 2)
	acc = collector.Accumulator(acc, 8)
	acc = collector.Accumulator(acc, 1)
	acc = collector.Accumulator(acc, 9)
	result := collector.Finisher(acc)
	if result.IsAbsent() || result.Get() != 9 {
		t.Errorf("expected 9, got %v", result)
	}
}

func TestGroupingBy(t *testing.T) {
	collector := GroupingBy(func(s string) rune { return rune(s[0]) })
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, "apple")
	acc = collector.Accumulator(acc, "apricot")
	acc = collector.Accumulator(acc, "banana")
	acc = collector.Accumulator(acc, "berry")
	acc = collector.Accumulator(acc, "cherry")
	result := collector.Finisher(acc)
	if len(result['a']) != 2 {
		t.Errorf("expected 2 words starting with 'a', got %d", len(result['a']))
	}
}

func TestPartitioningBy(t *testing.T) {
	collector := PartitioningBy(func(x int) bool { return x%2 == 0 })
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, 1)
	acc = collector.Accumulator(acc, 2)
	acc = collector.Accumulator(acc, 3)
	acc = collector.Accumulator(acc, 4)
	acc = collector.Accumulator(acc, 5)
	acc = collector.Accumulator(acc, 6)
	result := collector.Finisher(acc)
	if len(result[true]) != 3 {
		t.Errorf("expected 3 evens, got %d", len(result[true]))
	}
	if len(result[false]) != 3 {
		t.Errorf("expected 3 odds, got %d", len(result[false]))
	}
}

func TestToMap(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	collector := ToMap(
		func(p Person) string { return p.Name },
		func(p Person) int { return p.Age },
	)
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, Person{"Alice", 30})
	acc = collector.Accumulator(acc, Person{"Bob", 25})
	acc = collector.Accumulator(acc, Person{"Charlie", 35})
	result := collector.Finisher(acc)
	if result["Alice"] != 30 {
		t.Errorf("expected Alice's age to be 30, got %d", result["Alice"])
	}
}

func TestToMapWith(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	collector := ToMapWith(
		func(p Person) string { return p.Name },
		func(p Person) int { return p.Age },
		func(a, b int) int {
			if a > b {
				return a
			}
			return b
		},
	)
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, Person{"Alice", 30})
	acc = collector.Accumulator(acc, Person{"Alice", 25})
	acc = collector.Accumulator(acc, Person{"Bob", 35})
	result := collector.Finisher(acc)
	if result["Alice"] != 30 {
		t.Errorf("expected Alice's age to be 30 (max), got %d", result["Alice"])
	}
}

func TestSummarizing(t *testing.T) {
	collector := Summarizing(func(x int) float64 { return float64(x) })
	acc := collector.Supplier()
	acc = collector.Accumulator(acc, 1)
	acc = collector.Accumulator(acc, 2)
	acc = collector.Accumulator(acc, 3)
	acc = collector.Accumulator(acc, 4)
	acc = collector.Accumulator(acc, 5)
	result := collector.Finisher(acc)
	if result.Count != 5 {
		t.Errorf("expected count 5, got %d", result.Count)
	}
	if result.Sum != 15.0 {
		t.Errorf("expected sum 15, got %f", result.Sum)
	}
	if result.Min != 1.0 {
		t.Errorf("expected min 1, got %f", result.Min)
	}
	if result.Max != 5.0 {
		t.Errorf("expected max 5, got %f", result.Max)
	}
	if result.Average != 3.0 {
		t.Errorf("expected average 3, got %f", result.Average)
	}
}

func BenchmarkToSlice(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector := ToSlice[int]()
		acc := collector.Supplier()
		for _, v := range data {
			acc = collector.Accumulator(acc, v)
		}
		_ = collector.Finisher(acc)
	}
}

func BenchmarkToSet(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i % 100
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector := ToSet[int]()
		acc := collector.Supplier()
		for _, v := range data {
			acc = collector.Accumulator(acc, v)
		}
		_ = collector.Finisher(acc)
	}
}

func BenchmarkJoining(b *testing.B) {
	data := make([]string, 1000)
	for i := range data {
		data[i] = "word"
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector := Joining(", ")
		acc := collector.Supplier()
		for _, v := range data {
			acc = collector.Accumulator(acc, v)
		}
		_ = collector.Finisher(acc)
	}
}

func BenchmarkGroupingBy(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector := GroupingBy(func(x int) int { return x % 10 })
		acc := collector.Supplier()
		for _, v := range data {
			acc = collector.Accumulator(acc, v)
		}
		_ = collector.Finisher(acc)
	}
}

func BenchmarkSummarizing(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector := Summarizing(func(x int) float64 { return float64(x) })
		acc := collector.Supplier()
		for _, v := range data {
			acc = collector.Accumulator(acc, v)
		}
		_ = collector.Finisher(acc)
	}
}
