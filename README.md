# stdx

[![Go Reference](https://pkg.go.dev/badge/github.com/marouanesouiri/stdx.svg)](https://pkg.go.dev/github.com/marouanesouiri/stdx)
[![Go Report Card](https://goreportcard.com/badge/github.com/marouanesouiri/stdx)](https://goreportcard.com/report/github.com/marouanesouiri/stdx)
[![License: BSD 3-Clause](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

**stdx** is an extended standard library for Go. It adds data structures and tools that are common in other languages but missing in Go.

## 📦 Features

### Functional Tools
- **`stream`**: Process lists of data in a chain (filter, map, reduce), similar to Java Streams.
- **`optional`**: A safe way to handle values that might be missing, instead of using `nil`.
- **`either`**: A type that holds one of two possible values (usually a result or an error).
- **`tuple`**: Holds a fixed group of values (from 2 to 5 values) together.
- **`lazy`**: computes a value only when it is first needed, then remembers it.

### Data Structures
- **`cmap`**: A map that is safe to use from multiple parts of your code at the same time.
- **`omap`**: A map that remembers the order you added items.
- **`mmap`**: A map where one key can hold multiple values.
- **`set`**: A collection of unique items.

### Helpers
- **`scheduler`**: Runs tasks after a set delay using a single background worker.
- **`collectors`**: Helpers to convert Streams back into lists, maps, or sets.
- **`xlog`**: A simple, fast logger that supports JSON and text output.
- **`result`**: A way to handle success or failure without returning two values.

## 🚀 Installation

```bash
go get github.com/marouanesouiri/stdx
```

## 💡 Quick Start

### Stream API
```go
import (
    "fmt"
    "github.com/marouanesouiri/stdx/stream"
)

func main() {
    // Process a list of numbers
    nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    sum := stream.From(nums).
        Filter(func(n int) bool { return n%2 == 0 }). // Keep even numbers
        Map(func(n int) int { return n * 2 }).        // Double them
        Reduce(0, func(a, b int) int { return a + b }) // Add them up

    fmt.Println(sum) // Output: 60
}
```

### Optional Helper
```go
import (
    "fmt"
    "github.com/marouanesouiri/stdx/optional"
)

func FindUser(id int) optional.Option[string] {
    if id == 42 {
        return optional.Some("Alice")
    }
    return optional.None[string]()
}

func main() {
    user := FindUser(42)
    
    // Check if value exists
    if user.IsPresent() {
        fmt.Println("Found:", user.Get())
    }

    // Or provide a default
    name := user.OrElse("Unknown")
    fmt.Println("User:", name)
}
```

### Generic Set
```go
import "github.com/marouanesouiri/stdx/set"

s := set.New[string]()
s.Add("go")
s.Add("rust")
s.Add("go") // Duplicate is ignored

fmt.Println(s.Contains("go")) // true
fmt.Println(s.Size())         // 2
```

## 🤝 Contributing

Contributions are welcome!

## 📄 License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for deails.
