// Package lazy provides a generic Lazy type for delayed computation and memoization.
//
// A Lazy value represents a computation that is deferred until the value is actually needed.
// Once computed, the value is cached (memoized) and subsequent accesses return the cached value
// without re-executing the computation. All operations are thread-safe.
//
// # Basic Usage
//
// Create a Lazy value that computes on first access:
//
//	expensive := lazy.New(func() int {
//	    fmt.Println("Computing...")
//	    time.Sleep(1 * time.Second)
//	    return 42
//	})
//
//	// Nothing computed yet
//	fmt.Println("Created lazy value")
//
//	// Now it computes
//	val := expensive.Get() // prints "Computing..." and waits 1 second
//	fmt.Println(val)       // 42
//
//	// Uses cached value
//	val = expensive.Get() // returns immediately, no "Computing..." printed
//	fmt.Println(val)      // 42
//
// Create a Lazy value that's already computed:
//
//	immediate := lazy.Of(100)
//	val := immediate.Get() // returns immediately, no computation
//	fmt.Println(val)       // 100
//
// # Thread Safety
//
// Lazy values are safe for concurrent access. Multiple goroutines can call Get()
// simultaneously, and the supplier function will only execute once:
//
//	expensive := lazy.New(func() string {
//	    time.Sleep(100 * time.Millisecond)
//	    return "computed value"
//	})
//
//	// Multiple goroutines access concurrently
//	var wg sync.WaitGroup
//	for i := 0; i < 10; i++ {
//	    wg.Add(1)
//	    go func() {
//	        defer wg.Done()
//	        val := expensive.Get() // Only computes once total
//	        fmt.Println(val)
//	    }()
//	}
//	wg.Wait()
//
// # Checking Computation Status
//
// Check if a value has been computed without forcing computation:
//
//	lz := lazy.New(func() int { return 42 })
//
//	if !lz.IsComputed() {
//	    fmt.Println("Not computed yet")
//	}
//
//	val := lz.Get()
//
//	if lz.IsComputed() {
//	    fmt.Println("Now computed")
//	}
//
// # Transformations
//
// Transform lazy values without forcing computation:
//
//	lz := lazy.New(func() int {
//	    fmt.Println("Computing original")
//	    return 5
//	})
//
//	// Create transformed lazy value (no computation yet)
//	doubled := lazy.Map(lz, func(x int) int {
//	    fmt.Println("Doubling")
//	    return x * 2
//	})
//
//	// Nothing computed yet
//	fmt.Println("Transformations created")
//
//	// Now both compute
//	val := doubled.Get()
//	// prints: "Computing original"
//	// prints: "Doubling"
//	fmt.Println(val) // 10
//
// # Chaining Lazy Computations
//
// Use FlatMap to chain lazy computations:
//
//	config := lazy.New(func() string {
//	    fmt.Println("Loading config...")
//	    return "database.yml"
//	})
//
//	connection := lazy.FlatMap(config, func(configFile string) *lazy.Lazy[string] {
//	    return lazy.New(func() string {
//	        fmt.Println("Connecting to database from", configFile)
//	        return "connected"
//	    })
//	})
//
//	// Nothing computed yet
//	status := connection.Get()
//	// prints: "Loading config..."
//	// prints: "Connecting to database from database.yml"
//	fmt.Println(status) // "connected"
//
// # Filtering
//
// Apply predicates lazily:
//
//	lz := lazy.New(func() int { return 42 })
//
//	positive := lazy.Filter(lz, func(x int) bool {
//	    return x > 0
//	})
//
//	val := positive.Get() // 42 (predicate is true)
//
//	negative := lazy.Filter(lz, func(x int) bool {
//	    return x < 0
//	})
//
//	val = negative.Get() // 0 (zero value, predicate is false)
//

// # Use Cases
//
// Lazy is useful for:
//
// **Expensive Computations:**
//
//	// Only compute if needed
//	cachedResult := lazy.New(func() []byte {
//	    data, _ := os.ReadFile("large-file.dat")
//	    return data
//	})
//
//	if someCondition {
//	    data := cachedResult.Get() // File only read if condition is true
//	    process(data)
//	}
//
// **Database Connections:**
//
//	type App struct {
//	    db *lazy.Lazy[*sql.DB]
//	}
//
//	func NewApp() *App {
//	    return &App{
//	        db: lazy.New(func() *sql.DB {
//	            db, _ := sql.Open("postgres", "connection-string")
//	            return db
//	        }),
//	    }
//	}
//
//	func (a *App) Query() {
//	    db := a.db.Get() // Connection only opened when first query runs
//	    db.Query("SELECT ...")
//	}
//
// **Configuration Loading:**
//
//	var config = lazy.New(func() Config {
//	    fmt.Println("Loading configuration...")
//	    data, _ := os.ReadFile("config.json")
//	    var cfg Config
//	    json.Unmarshal(data, &cfg)
//	    return cfg
//	})
//
//	func GetSetting() string {
//	    return config.Get().Setting
//	}
//
// **Circular Dependencies:**
//
//	type ServiceA struct {
//	    b *lazy.Lazy[*ServiceB]
//	}
//
//	type ServiceB struct {
//	    a *lazy.Lazy[*ServiceA]
//	}
//
//	// Resolve circular dependency with lazy initialization
//	var a ServiceA
//	var b ServiceB
//	a.b = lazy.New(func() *ServiceB { return &b })
//	b.a = lazy.New(func() *ServiceA { return &a })
//
// **Memoization:**
//
//	func fibonacci(n int) *lazy.Lazy[int] {
//	    if n <= 1 {
//	        return lazy.Of(n)
//	    }
//	    return lazy.New(func() int {
//	        return fibonacci(n-1).Get() + fibonacci(n-2).Get()
//	    })
//	}
//
// # Performance Benefits
//
// **Deferred Execution:**
//   - Computations only run when needed
//   - Avoid unnecessary work in conditional code paths
//   - Faster startup times by deferring initialization
//
// **Memoization:**
//   - Expensive operations cached automatically
//   - Subsequent accesses are O(1)
//   - Memory for results only allocated when computed
//
// **Thread Safety Without Overhead:**
//   - Uses sync.Once for minimal synchronization cost
//   - No locks needed after first computation
//   - Multiple readers don't contend after initialization
//
// # Best Practices
//
// **Use Lazy for expensive operations:**
//
//	// Good - expensive operation
//	lazyData := lazy.New(func() []byte {
//	    return loadLargeFile()
//	})
//
//	// Bad - cheap operation (unnecessary overhead)
//	lazyInt := lazy.New(func() int {
//	    return 5 + 3
//	})
//
// **Thread-safe singletons:**
//
//	var instance = lazy.New(func() *Service {
//	    return &Service{
//	        // expensive initialization
//	    }
//	})
//
//	func GetService() *Service {
//	    return instance.Get() // Thread-safe singleton
//	}
//
// **Avoid side effects in suppliers:**
//
//	// Bad - side effects can happen at unexpected times
//	lz := lazy.New(func() int {
//	    fmt.Println("Side effect!") // This might not run when you expect
//	    return 42
//	})
//
//	// Good - pure computation
//	lz := lazy.New(func() int {
//	    return computeValue()
//	})
//
// # Comparison with Other Patterns
//
// **Lazy vs Regular Values:**
//
//	// Regular - always computed immediately
//	val := expensiveComputation()
//
//	// Lazy - computed only when needed
//	val := lazy.New(expensiveComputation)
//
// **Lazy vs sync.Once:**
//
//	// With sync.Once (manual)
//	var once sync.Once
//	var value int
//	once.Do(func() {
//	    value = compute()
//	})
//
//	// With Lazy (automatic)
//	value := lazy.New(compute)
//	val := value.Get()
//
// **Lazy vs Caching:**
//
//	// Manual caching
//	type Cache struct {
//	    mu    sync.Mutex
//	    value *int
//	}
//	func (c *Cache) Get() int {
//	    c.mu.Lock()
//	    defer c.mu.Unlock()
//	    if c.value == nil {
//	        v := compute()
//	        c.value = &v
//	    }
//	    return *c.value
//	}
//
//	// Lazy (simpler)
//	cache := lazy.New(compute)
//	val := cache.Get()
//
// # Memory Considerations
//
// Lazy values hold both the supplier function and the computed value in memory.
// After computation:
//   - The result is cached (uses memory)
//   - The supplier function reference remains (small overhead)
//
// For very large values, consider using pointers:
//
//	lazyLarge := lazy.New(func() *LargeStruct {
//	    return computeLargeStruct()
//	})
package lazy
