// Package tuple provides generic tuple types for grouping multiple values together.
//
// Tuples are fixed-size collections of values that may have different types.
// They provide a convenient way to group related values without defining a custom struct,
// making them particularly useful for temporary data structures, function return values,
// and working with paired or grouped data.
//
// # Available Tuple Types
//
// The package provides tuple types of various sizes:
//   - Tuple2[T, U] - A pair of two values (also called a pair)
//   - Tuple3[T, U, V] - A triple of three values
//   - Tuple4[T, U, V, W] - A quadruple of four values
//   - Tuple5[T, U, V, W, X] - A quintuple of five values
//
// # Basic Usage
//
// Create a Tuple2 (pair) using NewTuple2:
//
//	pair := tuple.NewTuple2("Alice", 30)
//	fmt.Println(pair.First)  // "Alice"
//	fmt.Println(pair.Second) // 30
//
// Create a Tuple3 (triple) using NewTuple3:
//
//	triple := tuple.NewTuple3("Alice", 30, "alice@example.com")
//	fmt.Println(triple.First)  // "Alice"
//	fmt.Println(triple.Second) // 30
//	fmt.Println(triple.Third)  // "alice@example.com"
//
// # Accessing Values
//
// Access tuple values directly through fields or use the Values method:
//
//	pair := tuple.NewTuple2("hello", 42)
//
//	// Direct field access
//	str := pair.First  // "hello"
//	num := pair.Second // 42
//
//	// Use Values() to get all values at once
//	str, num := pair.Values()
//
// # Transformations
//
// Transform individual elements or entire tuples using Map functions:
//
//	pair := tuple.NewTuple2(5, "hello")
//
//	// Transform the first element
//	doubled := pair.MapFirst(func(x int) int { return x * 2 })
//	fmt.Println(doubled.First) // 10
//
//	// Transform the second element
//	upper := pair.MapSecond(func(s string) string {
//	    return strings.ToUpper(s)
//	})
//	fmt.Println(upper.Second) // "HELLO"
//
//	// Transform both elements together
//	swapped := pair.Map(func(x int, s string) (int, string) {
//	    return x + len(s), s + "!"
//	})
//	fmt.Println(swapped.First, swapped.Second) // 10 "hello!"
//
//	// Map to different types
//	result := tuple.MapBoth(pair,
//	    func(x int) string { return fmt.Sprintf("%d", x) },
//	    func(s string) int { return len(s) },
//	)
//	fmt.Println(result.First, result.Second) // "5" 5
//
// # Swapping
//
// Swap the first and second elements of a Tuple2:
//
//	pair := tuple.NewTuple2("hello", 42)
//	swapped := pair.Swap()
//	fmt.Println(swapped.First, swapped.Second) // 42 "hello"
//
// # Comparison
//
// Compare tuples using custom equality functions:
//
//	pair1 := tuple.NewTuple2(5, "hello")
//	pair2 := tuple.NewTuple2(5, "hello")
//	pair3 := tuple.NewTuple2(10, "world")
//
//	equal := pair1.Equal(pair2,
//	    func(a, b int) bool { return a == b },
//	    func(a, b string) bool { return a == b },
//	)
//	fmt.Println(equal) // true
//
//	equal = pair1.Equal(pair3,
//	    func(a, b int) bool { return a == b },
//	    func(a, b string) bool { return a == b },
//	)
//	fmt.Println(equal) // false
//
// # Working with Slices
//
// Zip and unzip slices using tuple utilities:
//
//	names := []string{"Alice", "Bob", "Charlie"}
//	ages := []int{30, 25, 35}
//
//	// Zip two slices into a slice of tuples
//	people := tuple.Zip(names, ages)
//	for _, person := range people {
//	    fmt.Printf("%s is %d years old\n", person.First, person.Second)
//	}
//
//	// Unzip a slice of tuples back into two slices
//	names2, ages2 := tuple.Unzip(people)
//	fmt.Println(names2) // ["Alice", "Bob", "Charlie"]
//	fmt.Println(ages2)  // [30, 25, 35]
//
//	// Zip three slices
//	emails := []string{"alice@example.com", "bob@example.com", "charlie@example.com"}
//	users := tuple.Zip3(names, ages, emails)
//	for _, user := range users {
//	    fmt.Printf("%s (%d) - %s\n", user.First, user.Second, user.Third)
//	}
//
//	// Unzip three slices
//	names3, ages3, emails3 := tuple.Unzip3(users)
//
// # JSON Serialization
//
// Tuples automatically support JSON marshaling and unmarshaling.
// They are serialized as JSON arrays:
//
//	type Person struct {
//	    Name string
//	    Info tuple.Tuple2[int, string] // age and email
//	}
//
//	person := Person{
//	    Name: "Alice",
//	    Info: tuple.NewTuple2(30, "alice@example.com"),
//	}
//
//	data, _ := json.Marshal(person)
//	// {"Name":"Alice","Info":[30,"alice@example.com"]}
//
//	var decoded Person
//	json.Unmarshal(data, &decoded)
//	fmt.Println(decoded.Info.First)  // 30
//	fmt.Println(decoded.Info.Second) // "alice@example.com"
//
// # String Representation
//
// All tuple types implement the String method for easy printing:
//
//	pair := tuple.NewTuple2("hello", 42)
//	fmt.Println(pair) // (hello, 42)
//
//	triple := tuple.NewTuple3(1, 2, 3)
//	fmt.Println(triple) // (1, 2, 3)
//
// # When to Use Tuples
//
// Tuples are most useful when:
//   - Returning multiple values from a function without defining a custom struct
//   - Temporarily grouping related values during data processing
//   - Working with paired or grouped data (like key-value pairs)
//   - Zipping and unzipping slices for parallel iteration
//   - Representing coordinates, ranges, or intervals
//   - Building more complex data structures
//
// Example use cases:
//   - Coordinate pairs: tuple.NewTuple2(x, y)
//   - RGB colors: tuple.NewTuple3(red, green, blue)
//   - Function results: tuple.NewTuple2(result, error)
//   - Dictionary entries: tuple.NewTuple2(key, value)
//   - Ranges: tuple.NewTuple2(min, max)
//
// # Common Patterns
//
// Returning multiple heterogeneous values:
//
//	func parseUser(data string) tuple.Tuple2[string, int] {
//	    // Parse user data...
//	    return tuple.NewTuple2(name, age)
//	}
//
//	user := parseUser(data)
//	fmt.Printf("Name: %s, Age: %d\n", user.First, user.Second)
//
// Iterating over paired data:
//
//	keys := []string{"a", "b", "c"}
//	values := []int{1, 2, 3}
//
//	for _, pair := range tuple.Zip(keys, values) {
//	    fmt.Printf("%s = %d\n", pair.First, pair.Second)
//	}
//
// Building lookup tables:
//
//	type Config struct {
//	    Settings []tuple.Tuple2[string, string]
//	}
//
//	config := Config{
//	    Settings: []tuple.Tuple2[string, string]{
//	        tuple.NewTuple2("host", "localhost"),
//	        tuple.NewTuple2("port", "8080"),
//	        tuple.NewTuple2("debug", "true"),
//	    },
//	}
//
// # Performance Considerations
//
// Tuples are value types (structs) and are passed by value. For large tuples or
// tuples containing large values, consider using pointers if performance is critical.
//
// The Zip and Unzip functions allocate new slices. For performance-critical code
// with large datasets, consider using index-based iteration instead.
//
// # Type Safety
//
// Tuples provide compile-time type safety. Each element's type is known at compile time,
// making them safer than using []interface{} or map[string]interface{} for grouping
// heterogeneous data.
//
//	// Compile-time type checking
//	pair := tuple.NewTuple2("Alice", 30)
//	name := pair.First  // Type is string
//	age := pair.Second  // Type is int
//
// # Comparison with Other Approaches
//
// Instead of using a custom struct:
//
//	// Without tuples
//	type Result struct {
//	    Value int
//	    Error error
//	}
//	func process() Result { /* ... */ }
//
//	// With tuples
//	func process() tuple.Tuple2[int, error] { /* ... */ }
//
// Instead of using multiple return values:
//
//	// Multiple returns (limited to same usage context)
//	func getCoords() (int, int) {
//	    return 10, 20
//	}
//	x, y := getCoords()
//
//	// With tuples (can be stored, passed around)
//	func getCoords() tuple.Tuple2[int, int] {
//	    return tuple.NewTuple2(10, 20)
//	}
//	coords := getCoords()
//	// Later...
//	x, y := coords.Values()
package tuple
