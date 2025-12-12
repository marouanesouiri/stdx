// Package optional provides a generic Option type for representing optional values.
//
// The Option type is a container that either holds a value (Some) or is empty (None).
// This provides a type-safe alternative to using pointers or special "zero" sentinel values
// to represent the absence of a value, helping to avoid nil pointer dereferences and
// making the intent of your code more explicit.
//
// # Basic Usage
//
// Create an Option with a value using Some:
//
//	opt := optional.Some(42)
//	fmt.Println(opt.IsPresent()) // true
//	fmt.Println(opt.Get())       // 42
//
// Create an empty Option using None:
//
//	opt := optional.None[string]()
//	fmt.Println(opt.IsPresent()) // false
//	fmt.Println(opt.OrEmpty())   // "" (zero value)
//
// # Creating Options
//
// Several constructors are available for creating Options:
//   - Some(value) - creates an Option containing the given value
//   - None[T]() - creates an empty Option
//   - FromPtr(ptr) - creates an Option from a pointer (nil becomes None)
//   - FromZero(value) - creates an Option, treating zero values as None
//
// Example:
//
//	// From pointer
//	var ptr *int = nil
//	opt1 := optional.FromPtr(ptr) // None
//
//	val := 42
//	opt2 := optional.FromPtr(&val) // Some(42)
//
//	// From zero value
//	opt3 := optional.FromZero(0)  // None
//	opt4 := optional.FromZero(42) // Some(42)
//
// # Checking Presence
//
// Use IsPresent() or IsAbsent() to check if a value exists:
//
//	opt := optional.Some(42)
//	if opt.IsPresent() {
//	    fmt.Println("Value:", opt.Get())
//	}
//
// # Retrieving Values
//
// Several methods exist for retrieving values safely:
//   - Get() - returns the value (even if absent, returns zero value)
//   - MustGet() - returns the value or panics if absent
//   - GetErr() - returns the value and an error if absent
//   - OrElse(fallback) - returns the value or a fallback
//   - OrElseGet(supplier) - returns the value or calls a supplier function
//   - OrElseErr(supplier) - returns the value or an error from a supplier
//   - OrEmpty() - returns the value or the zero value of T
//   - OrPanic(message) - returns the value or panics with a custom message
//   - Ptr() - returns a pointer to the value, or nil if absent
//
// Example:
//
//	opt := optional.None[int]()
//
//	// Safe retrieval with fallback
//	val1 := opt.OrElse(42)                    // 42
//	val2 := opt.OrElseGet(func() int { return 100 }) // 100
//
//	// Error handling
//	val3, err := opt.OrElseErr(func() error {
//	    return errors.New("value not found")
//	})
//	if err != nil {
//	    fmt.Println("Error:", err)
//	}
//
// # Conditional Execution
//
// Execute code conditionally based on the presence of a value:
//
//	opt := optional.Some(42)
//
//	// Execute if present
//	opt.IfPresent(func(val int) {
//	    fmt.Println("Value:", val)
//	})
//
//	// Execute different code based on presence
//	opt.IfPresentOrElse(
//	    func(val int) { fmt.Println("Found:", val) },
//	    func() { fmt.Println("Not found") },
//	)
//
// # Transformations
//
// Transform the value inside an Option using functional methods:
//
//	opt := optional.Some(5)
//
//	// Map transforms the value
//	doubled := optional.Map(opt, func(x int) int { return x * 2 })
//	fmt.Println(doubled.Get()) // 10
//
//	// FlatMap chains optional-returning functions
//	result := optional.FlatMap(opt, func(x int) optional.Option[int] {
//	    if x > 0 {
//	        return optional.Some(x * 10)
//	    }
//	    return optional.None[int]()
//	})
//	fmt.Println(result.Get()) // 50
//
//	// TryMap handles operations that may fail
//	converted := optional.TryMap(optional.Some("42"), func(s string) (int, error) {
//	    return strconv.Atoi(s)
//	})
//	fmt.Println(converted.Get()) // 42
//
// # Filtering
//
// Filter values based on a predicate:
//
//	opt := optional.Some(42)
//	positive := opt.Filter(func(x int) bool { return x > 0 })
//	fmt.Println(positive.IsPresent()) // true
//
//	negative := opt.Filter(func(x int) bool { return x < 0 })
//	fmt.Println(negative.IsPresent()) // false
//
// # Combining Options
//
// Combine multiple Options using logical operations:
//
//	opt1 := optional.Some(1)
//	opt2 := optional.Some(2)
//	none := optional.None[int]()
//
//	// And: returns second if first is Some, otherwise None
//	result1 := opt1.And(opt2) // Some(2)
//	result2 := none.And(opt1) // None
//
//	// Or: returns first if Some, otherwise second
//	result3 := opt1.Or(opt2)  // Some(1)
//	result4 := none.Or(opt1)  // Some(1)
//
//	// Xor: Some if exactly one is Some
//	result5 := opt1.Xor(none) // Some(1)
//	result6 := opt1.Xor(opt2) // None (both are Some)
//

// # JSON Serialization
//
// Options automatically support JSON marshaling and unmarshaling:
//
//	type User struct {
//	    Name     string
//	    Email    optional.Option[string]
//	    Age      optional.Option[int]
//	}
//
//	user := User{
//	    Name:  "Alice",
//	    Email: optional.Some("alice@example.com"),
//	    Age:   optional.None[int](),
//	}
//
//	data, _ := json.Marshal(user)
//	// {"Name":"Alice","Email":"alice@example.com","Age":null}
//
//	var decoded User
//	json.Unmarshal(data, &decoded)
//	fmt.Println(decoded.Email.OrEmpty()) // "alice@example.com"
//	fmt.Println(decoded.Age.IsPresent()) // false
//
// # Comparison
//
// Compare two Options using a custom equality function:
//
//	opt1 := optional.Some(42)
//	opt2 := optional.Some(42)
//	opt3 := optional.Some(100)
//
//	eq := func(a, b int) bool { return a == b }
//	fmt.Println(opt1.Equal(opt2, eq)) // true
//	fmt.Println(opt1.Equal(opt3, eq)) // false
//
// # When to Use
//
// Use Option when:
//   - A function may or may not return a value (instead of returning nil or a special sentinel)
//   - You want to make the absence of a value explicit in your type signature
//   - You want to avoid nil pointer dereferences
//   - You want to chain operations that may fail without excessive error checking
//   - You need optional fields in structs that serialize cleanly to JSON
//
// Example use cases:
//   - Configuration values that may be unset
//   - Cache lookups that may miss
//   - User input that may be empty
//   - API responses with optional fields
//   - Search results that may not find anything
package optional
