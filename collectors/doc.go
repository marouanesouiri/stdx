// Package collectors provides reusable Collector implementations for the stream package.
//
// Collectors are used with stream.CollectWith() to perform common reduction operations
// on streams in a composable way, similar to Java's Collectors utility class.
//
// # Basic Usage
//
// Use collectors with stream.CollectWith():
//
//	import (
//	    "github.com/marouanesouiri/stdx/stream"
//	    "github.com/marouanesouiri/stdx/collectors"
//	)
//
//	result := stream.From(data).
//	    Filter(predicate).
//	    CollectWith(collectors.ToSlice[int]())
//
// # Available Collectors
//
// Collection Collectors:
//   - ToSlice: Collect elements into a slice
//   - ToSet: Collect elements into a Set (removes duplicates)
//
// String Collectors:
//   - Joining: Join strings with a separator
//   - JoiningWith: Join strings with separator, prefix, and suffix
//
// Numeric Collectors:
//   - Counting: Count the number of elements
//   - Summing: Sum numeric values extracted by a mapper function
//   - Averaging: Compute the average of numeric values
//   - MinBy: Find the minimum element according to a comparator
//   - MaxBy: Find the maximum element according to a comparator
//
// Grouping Collectors:
//   - GroupingBy: Group elements by a key function
//   - PartitioningBy: Partition elements into two groups based on a predicate
//   - ToMap: Collect elements into a map
//   - ToMapWith: Collect into a map with a merge function for duplicate keys
//
// Statistical Collectors:
//   - Summarizing: Compute count, sum, min, max, and average in one pass
//
// # Examples
//
// Joining strings:
//
//	words := []string{"apple", "banana", "cherry"}
//	csv := stream.CollectWith(
//	    stream.From(words),
//	    collectors.Joining(", "),
//	)  // "apple, banana, cherry"
//
//	list := stream.CollectWith(
//	    stream.From(words),
//	    collectors.JoiningWith(", ", "[", "]"),
//	)  // "[apple, banana, cherry]"
//
// Grouping:
//
//	words := []string{"apple", "apricot", "banana", "berry"}
//	byFirstLetter := stream.CollectWith(
//	    stream.From(words),
//	    collectors.GroupingBy(func(s string) rune { return rune(s[0]) }),
//	)
//	// map[rune][]string{'a': ["apple", "apricot"], 'b': ["banana", "berry"]}
//
// Statistics:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	stats := stream.CollectWith(
//	    stream.From(numbers),
//	    collectors.Summarizing(func(x int) float64 { return float64(x) }),
//	)
//	// Statistics{Count: 5, Sum: 15, Min: 1, Max: 5, Average: 3}
//
// Partitioning:
//
//	numbers := []int{1, 2, 3, 4, 5, 6}
//	partitioned := stream.CollectWith(
//	    stream.From(numbers),
//	    collectors.PartitioningBy(func(x int) bool { return x%2 == 0 }),
//	)
//	// map[bool][]int{true: [2, 4, 6], false: [1, 3, 5]}
//
// Custom mapping:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	people := []Person{
//	    {"Alice", 30},
//	    {"Bob", 25},
//	    {"Charlie", 35},
//	}
//
//	ageByName := stream.CollectWith(
//	    stream.From(people),
//	    collectors.ToMap(
//	        func(p Person) string { return p.Name },
//	        func(p Person) int { return p.Age },
//	    ),
//	)
//	// map[string]int{"Alice": 30, "Bob": 25, "Charlie": 35}
//
// Handling duplicates with merge function:
//
//	people := []Person{
//	    {"Alice", 30},
//	    {"Alice", 25},  // duplicate name
//	    {"Bob", 35},
//	}
//
//	maxAgeByName := stream.CollectWith(
//	    stream.From(people),
//	    collectors.ToMapWith(
//	        func(p Person) string { return p.Name },
//	        func(p Person) int { return p.Age },
//	        func(age1, age2 int) int {
//	            if age1 > age2 {
//	                return age1
//	            }
//	            return age2
//	        },
//	    ),
//	)
//	// map[string]int{"Alice": 30, "Bob": 35}  // kept max age for Alice
//
// Numeric operations:
//
//	numbers := []int{1, 2, 3, 4, 5}
//
//	sum := stream.CollectWith(
//	    stream.From(numbers),
//	    collectors.Summing(func(x int) int { return x }),
//	)  // 15
//
//	avg := stream.CollectWith(
//	    stream.From(numbers),
//	    collectors.Averaging(func(x int) float64 { return float64(x) }),
//	)  // 3.0
//
//	min := stream.CollectWith(
//	    stream.From(numbers),
//	    collectors.MinBy(func(a, b int) bool { return a < b }),
//	)  // Some(1)
//
// # Performance
//
// Collectors are designed to be efficient:
//   - Single-pass processing (except for composite collectors)
//   - Pre-allocated data structures where possible
//   - Minimal allocations
//
// For the best performance, choose the most specific collector for your use case
// rather than composing multiple operations.
package collectors
