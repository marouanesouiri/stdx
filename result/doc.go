/*
Package result provides a type that represents either success (Ok) or failure (Err).

It is a safer alternative to returning (T, error) and allows for more expressive error handling.
It integrates seamlessly with standard Go code using From() and ToPair().

# Usage

Create a Result:

	r1 := result.Ok(42)
	r2 := result.Err[int](errors.New("oops"))

	// From standard Go return values
	r3 := result.From(os.Open("file.txt"))

Check status:

	if r.IsOk() {
		fmt.Println(r.Value())
	}

Safe usage:

	val := r.UnwrapOr(0)
	val := r.UnwrapOrElse(func() int { return calculateDefault() })

Interop:

	// Convert back to (T, error)
	val, err := r.ToPair()

	// Convert to Option
	opt := r.Option() // Some(val) or None
*/
package result
