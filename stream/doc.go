// Package stream provides a functional Stream API for Go, similar to Java's Stream API.
//
// Streams enable functional-style operations on sequences of elements, supporting both
// intermediate operations (filter, map, etc.) and terminal operations (collect, reduce, etc.).
// All operations use lazy evaluation, processing elements only when a terminal operation is invoked.
//
// # Go 1.23 Integration
//
// This package leverages Go 1.23's iter package for standardized iteration.
// Streams wrap iter.Seq[T] and can be used directly with for-range loops:
//
//	s := stream.From([]int{1, 2, 3, 4, 5})
//	for v := range s.Filter(func(x int) bool { return x > 2 }).Seq() {
//	    fmt.Println(v) // 3, 4, 5
//	}
//
// # Creating Streams
//
// Multiple ways to create streams:
//
//	// From slice
//	s1 := stream.From([]int{1, 2, 3, 4, 5})
//
//	// From variadic args
//	s2 := stream.Of(1, 2, 3, 4, 5)
//
//	// From range
//	s3 := stream.Range(1, 10)
//
//	// From channel
//	ch := make(chan int)
//	s4 := stream.FromChannel(ch)
//
//	// From iter.Seq
//	s5 := stream.FromSeq(someIterSeq)
//
//	// Empty stream
//	s6 := stream.Empty[int]()
//
//	// Infinite streams
//	s7 := stream.Generate(func() int { return rand.Intn(100) })
//	s8 := stream.Iterate(0, func(n int) int { return n + 1 })
//
// # Intermediate Operations
//
// Intermediate operations return a new Stream and are lazily evaluated:
//
//	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//	result := stream.From(data).
//	    Filter(func(x int) bool { return x%2 == 0 }).  // Keep evens
//	    Limit(3).                                       // Take first 3
//	    ToSlice()                                       // [2, 4, 6]
//
// Available intermediate operations:
//   - Filter: Keep elements matching predicate
//   - Map: Transform each element
//   - FlatMap: Transform and flatten nested streams
//   - Distinct: Remove duplicates
//   - DistinctBy: Remove duplicates by key function
//   - Sorted: Sort elements
//   - Peek: Perform action without modification
//   - Limit: Take first n elements
//   - Skip: Skip first n elements
//   - TakeWhile: Take while predicate is true
//   - DropWhile: Drop while predicate is true
//   - Concat: Concatenate with another stream
//   - Reverse: Reverse element order
//
// # Terminal Operations
//
// Terminal operations consume the stream and produce a result:
//
//	data := []int{1, 2, 3, 4, 5}
//
//	// Collect using collector
//	slice := stream.Collect(stream.From(data), collectors.ToSlice[int]())
//
//	// Or use ToSlice() directly
//	slice := stream.From(data).ToSlice()
//
//	// Reduce to single value
//	sum := stream.From(data).Reduce(0, func(a, b int) int { return a + b })
//
//	// Count elements
//	count := stream.From(data).Count()
//
//	// Find first matching
//	first := stream.From(data).
//	    Filter(func(x int) bool { return x > 3 }).
//	    FindFirst()  // Some(4)
//
//	// Check conditions
//	allEven := stream.From(data).AllMatch(func(x int) bool { return x%2 == 0 })
//	hasEven := stream.From(data).AnyMatch(func(x int) bool { return x%2 == 0 })
//
// # Using Collectors
//
// Collectors provide reusable reduction operations from the collectors package:
//
//	import "github.com/marouanesouiri/stdx/collectors"
//
//	// Join strings
//	csv := stream.Collect(
//	    stream.From([]string{"a", "b", "c"}),
//	    collectors.Joining(", "),
//	)  // "a, b, c"
//
//	// Group by key
//	grouped := stream.Collect(
//	    stream.From(words),
//	    collectors.GroupingBy(func(s string) int { return len(s) }),
//	)
//
//	// Collect statistics
//	stats := stream.Collect(
//	    stream.From(numbers),
//	    collectors.Summarizing(func(x int) float64 { return float64(x) }),
//	)

// # Type Transformations
//
// The Map operation can change element types:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	strings := stream.Map(
//	    stream.From(numbers),
//	    func(n int) string { return strconv.Itoa(n) },
//	).ToSlice()  // ["1", "2", "3", "4", "5"]
//
// # Lazy Evaluation
//
// Streams use lazy evaluation - intermediate operations build a pipeline
// without processing elements. Processing happens only when a terminal operation is called:
//
//	count := 0
//	s := stream.From([]int{1, 2, 3, 4, 5}).
//	    Peek(func(x int) { count++ }).
//	    Filter(func(x int) bool { return x%2 == 0 })
//
//	// count is still 0 - nothing executed yet
//
//	result := s.ToSlice()
//	// Now count is 5 (all elements processed)
//
// # Integration with Custom Types
//
// Custom types can implement Streamer or Collectable interfaces:
//
//	type MyList[T any] struct { /* ... */ }
//
//	// Implement Streamer to create streams from your type
//	func (l MyList[T]) Stream() stream.Stream[T] {
//	    return stream.FromSeq(func(yield func(T) bool) {
//	        // yield elements
//	    })
//	}
//
//	// Implement Collectable for smart collection
//	func (l MyList[T]) Collect(seq iter.Seq[T]) MyList[T] {
//	    result := NewMyList[T]()
//	    for v := range seq {
//	        result.Append(v)
//	    }
//	    return result
//	}
//
// # Performance Considerations
//
// Streams are designed for readability and functional composition. For maximum performance
// in hot paths, consider:
//   - Using direct loops for simple operations
//   - Avoiding unnecessary intermediate operations
//   - Using Limit() early to reduce processing
//   - Being aware that operations like Sorted() and Distinct() must materialize the entire stream
//
// # Examples
//
// Filter and transform:
//
//	type User struct {
//	    Name   string
//	    Active bool
//	    Age    int
//	}
//
//	activeEmails := stream.From(users).
//	    Filter(func(u User) bool { return u.Active }).
//	    Map(func(u User) string { return u.Email }).
//	    ToSlice()
//
// Complex pipeline:
//
//	result := stream.From(data).
//	    Filter(func(x int) bool { return x > 0 }).
//	    Distinct().
//	    Sorted(func(a, b int) bool { return a < b }).
//	    Limit(10).
//	    ToSlice()
//
// Working with grouping:
//
//	byAge := stream.GroupBy(
//	    stream.From(users),
//	    func(u User) int { return u.Age },
//	)  // map[int][]User
package stream
