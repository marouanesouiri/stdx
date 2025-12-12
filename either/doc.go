// Package either provides a generic Either type for representing values that can be one of two types.
//
// Either is a type-safe way to represent a disjoint union of two types. An Either value is
// either Left or Right, but never both. By convention, Left is typically used for error or
// failure cases, while Right is used for success cases, though both represent valid values.
//
// This is different from Option, which represents presence or absence. Either always contains
// a value - it's just a question of which type.
//
// # Basic Usage
//
// Create an Either with a left value (conventionally used for errors):
//
//	result := either.Left[error, int](errors.New("failed"))
//	fmt.Println(result.IsLeft())  // true
//	fmt.Println(result.IsRight()) // false
//
// Create an Either with a right value (conventionally used for success):
//
//	result := either.Right[error, int](42)
//	fmt.Println(result.IsLeft())  // false
//	fmt.Println(result.IsRight()) // true
//	fmt.Println(result.Right())   // 42
//
// # Left vs Right Convention
//
// The convention is:
//   - **Left** = error, failure, or alternative path
//   - **Right** = success, expected value, or primary path
//
// This mirrors how Result types work in Rust and other languages. However, both sides
// are equally valid - Unlike Option, there's no "absent" state.
//
//	// Parsing example
//	func parseInt(s string) either.Either[error, int] {
//	    val, err := strconv.Atoi(s)
//	    if err != nil {
//	        return either.Left[error, int](err)
//	    }
//	    return either.Right[error, int](val)
//	}
//
// # Checking Which Side
//
// Use IsLeft() or IsRight() to determine which value is present:
//
//	result := parseInt("42")
//	if result.IsRight() {
//	    fmt.Println("Success:", result.Right())
//	} else {
//	    fmt.Println("Error:", result.Left())
//	}
//
// # Retrieving Values Safely
//
// Several methods exist for safe value retrieval:
//
//	result := either.Left[string, int]("error")
//
//	// With fallback
//	val := result.RightOr(0)           // 0 (fallback)
//	err := result.LeftOr("no error")   // "error"
//
//	// With boolean check
//	if val, ok := result.GetRight(); ok {
//	    fmt.Println("Got value:", val)
//	}
//
//	if err, ok := result.GetLeft(); ok {
//	    fmt.Println("Got error:", err)
//	}
//
// # Pattern Matching with Fold
//
// Use Fold to handle both cases in one expression:
//
//	result := parseInt("42")
//	message := either.Fold(result,
//	    func(err error) string {
//	        return "Error: " + err.Error()
//	    },
//	    func(val int) string {
//	        return fmt.Sprintf("Success: %d", val)
//	    },
//	)
//	fmt.Println(message)
//
// # Transformations
//
// Transform the right value (success path) using Map:
//
//	result := either.Right[error, int](5)
//	doubled := result.Map(func(x int) int { return x * 2 })
//	fmt.Println(doubled.Right()) // 10
//
// Transform the left value (error path) using MapLeft:
//
//	result := either.Left[string, int]("failed")
//	wrapped := result.MapLeft(func(s string) string {
//	    return "Error: " + s
//	})
//	fmt.Println(wrapped.Left()) // "Error: failed"
//
// Transform to different types using MapBoth:
//
//	result := either.Right[error, int](42)
//	transformed := either.MapBoth(result,
//	    func(err error) string { return err.Error() },
//	    func(val int) string { return fmt.Sprintf("%d", val) },
//	)
//	// transformed is Either[string, string]
//
// # Chaining Operations
//
// Use FlatMap to chain operations that return Either:
//
//	result := either.Right[error, int](5)
//	final := either.FlatMap(result, func(x int) either.Either[error, int] {
//	    if x > 0 {
//	        return either.Right[error, int](x * 10)
//	    }
//	    return either.Left[error, int](errors.New("value must be positive"))
//	})
//	fmt.Println(final.Right()) // 50
//
// # Swapping Sides
//
// Swap left and right values:
//
//	result := either.Left[string, int]("error")
//	swapped := result.Swap()
//	// swapped is Either[int, string] with Right("error")
//	fmt.Println(swapped.IsRight()) // true
//	fmt.Println(swapped.Right())   // "error"
//
// # JSON Serialization
//
// Either values serialize to JSON with type information:
//
//	type Response struct {
//	    Result either.Either[string, int]
//	}
//
//	// Success case
//	resp := Response{Result: either.Right[string, int](42)}
//	data, _ := json.Marshal(resp)
//	// {"Result":{"type":"right","value":42}}
//
//	// Error case
//	resp = Response{Result: either.Left[string, int]("failed")}
//	data, _ = json.Marshal(resp)
//	// {"Result":{"type":"left","value":"failed"}}
//
//	// Unmarshal
//	var decoded Response
//	json.Unmarshal(data, &decoded)
//	if decoded.Result.IsLeft() {
//	    fmt.Println("Error:", decoded.Result.Left())
//	}
//
// # When to Use Either
//
// Use Either when:
//   - You need to return one of two distinct types
//   - Both outcomes are valid (not "present" vs "absent")
//   - You want type-safe error handling with custom error types
//   - You want to avoid interface{} or type assertions
//   - You need to represent branching logic in types
//
// Example use cases:
//   - Parsing functions that return parsed value or error description
//   - Validation that returns validated value or validation errors
//   - Operations with alternative outcomes (not just success/failure)
//   - Type-safe tagged unions
//   - Railway-oriented programming patterns
//
// # Common Patterns
//
// Error handling with Either:
//
//	func divide(a, b int) either.Either[error, int] {
//	    if b == 0 {
//	        return either.Left[error, int](errors.New("division by zero"))
//	    }
//	    return either.Right[error, int](a / b)
//	}
//
//	result := divide(10, 2)
//	value := result.RightOr(0)
//
// Validation with multiple error types:
//
//	type ValidationError struct {
//	    Field   string
//	    Message string
//	}
//
//	func validateAge(age int) either.Either[ValidationError, int] {
//	    if age < 0 || age > 150 {
//	        return either.Left[ValidationError, int](ValidationError{
//	            Field:   "age",
//	            Message: "age must be between 0 and 150",
//	        })
//	    }
//	    return either.Right[ValidationError, int](age)
//	}
//
// Chaining validations:
//
//	result := validateAge(25)
//	doubled := either.FlatMap(result, func(age int) either.Either[ValidationError, int] {
//	    if age < 18 {
//	        return either.Left[ValidationError, int](ValidationError{
//	            Field:   "age",
//	            Message: "must be 18 or older",
//	        })
//	    }
//	    return either.Right[ValidationError, int](age * 2)
//	})
//
// # Either vs Option vs Result
//
// **Option[T]** - Represents presence or absence of a value
//
//	opt := optional.Some(42)  // Value is present
//	opt := optional.None[int]() // Value is absent
//
// **Either[L, R]** - Represents one of two possible values
//
//	either := either.Right[error, int](42) // Right value
//	either := either.Left[error, int](err) // Left value
//
// **Result[T, E]** (if implemented) - Specialized Either for Ok/Err
//
//	result := result.Ok[int, error](42)   // Success
//	result := result.Err[int, error](err) // Failure
//
// Either is more general than Result - both sides can be any type,
// while Result assumes one side is an error type.
//
// # Performance Considerations
//
// Either is a value type (struct) containing both L and R fields. Memory usage
// includes space for both types, even though only one is valid at a time.
//
// For large types, consider using Either[*Large, *OtherLarge] with pointers.
//
// # Type Safety
//
// Either provides compile-time type safety for union types:
//
//	// Type-safe branching
//	var result either.Either[string, int]
//	if someCondition {
//	    result = either.Left[string, int]("error message")
//	} else {
//	    result = either.Right[string, int](42)
//	}
//
//	// Both branches are type-checked
//	message := either.Fold(result,
//	    func(err string) string { return err },
//	    func(val int) string { return fmt.Sprintf("%d", val) },
//	)
package either
