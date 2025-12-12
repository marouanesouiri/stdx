// Package omap provides an ordered map that maintains insertion order.
//
// OrderedMap is a map implementation that preserves the order in which keys were inserted.
// Unlike Go's built-in map, iteration over an OrderedMap is predictable and follows insertion order.
// This is similar to Java's LinkedHashMap.
//
// # Basic Usage
//
// Create and use an ordered map:
//
//	m := omap.New[string, int]()
//	m.Set("first", 1)
//	m.Set("second", 2)
//	m.Set("third", 3)
//
//	for _, key := range m.Keys() {
//	    fmt.Println(key) // Prints: first, second, third (in order)
//	}
//
// # Order Preservation
//
// The map maintains insertion order. Updating an existing key moves it to the end:
//
//	m := omap.New[string, int]()
//	m.Set("a", 1)
//	m.Set("b", 2)
//	m.Set("a", 10) // Updates value and moves to end
//
//	keys := m.Keys() // ["b", "a"]
//
// # Basic Operations
//
//	m := omap.New[string, int]()
//
//	m.Set("key", 42)
//	val, ok := m.Get("key") // 42, true
//	exists := m.Has("key")  // true
//	m.Delete("key")         // true
//	size := m.Len()
//
// # Iteration
//
// Iterate in insertion order:
//
//	m := omap.New[string, int]()
//	m.Set("z", 1)
//	m.Set("a", 2)
//	m.Set("m", 3)
//
//	m.Range(func(key string, value int) bool {
//	    fmt.Printf("%s=%d ", key, value)
//	    return true // Continue iteration
//	})
//
// Get ordered slices:
//
//	keys := m.Keys()     // ["z", "a", "m"] in insertion order
//	values := m.Values() // [1, 2, 3] in insertion order
//	items := m.Items()   // All key-value pairs in order
//
// # First and Last
//
// Access the first or last inserted entries:
//
//	m := omap.New[string, int]()
//	m.Set("first", 1)
//	m.Set("second", 2)
//	m.Set("third", 3)
//
//	key, val, ok := m.First() // "first", 1, true
//	key, val, ok = m.Last()   // "third", 3, true
//
// # Pop Operations
//
// Remove and return first or last entries:
//
//	m := omap.New[string, int]()
//	m.Set("a", 1)
//	m.Set("b", 2)
//	m.Set("c", 3)
//
//	key, val, ok := m.PopFirst() // "a", 1, true
//	key, val, ok = m.PopLast()   // "c", 3, true
//
// # Use Cases
//
// **LRU Cache:**
//
//	type LRUCache struct {
//	    capacity int
//	    cache    *omap.OrderedMap[string, any]
//	}
//
//	func (c *LRUCache) Get(key string) (any, bool) {
//	    val, ok := c.cache.Get(key)
//	    if ok {
//	        c.cache.Set(key, val) // Move to end
//	    }
//	    return val, ok
//	}
//
//	func (c *LRUCache) Put(key string, value any) {
//	    c.cache.Set(key, value)
//	    if c.cache.Len() > c.capacity {
//	        c.cache.PopFirst() // Evict oldest
//	    }
//	}
//
// **Ordered Configuration:**
//
//	config := omap.New[string, string]()
//	config.Set("host", "localhost")
//	config.Set("port", "8080")
//	config.Set("debug", "true")
//
//	for _, key := range config.Keys() {
//	    fmt.Printf("%s = %s\n", key, config.Get(key))
//	}
//
// **Maintaining History:**
//
//	history := omap.New[time.Time, Event]()
//	history.Set(time.Now(), event1)
//	history.Set(time.Now(), event2)
//
//	for _, item := range history.Items() {
//	    fmt.Printf("%v: %v\n", item.Key, item.Value)
//	}
//
// # Performance
//
// **Time Complexity:**
//   - Set: O(1)
//   - Get: O(1)
//   - Delete: O(1)
//   - First/Last: O(1)
//   - PopFirst/PopLast: O(1)
//   - Keys/Values/Items: O(n)
//   - Range: O(n)
//
// **Space Complexity:**
//   - O(n) where n is the number of entries
//   - Additional overhead per entry: 3 pointers (prev, next, map pointer)
//
// # Comparison with map
//
//	// Regular map - undefined iteration order
//	m := make(map[string]int)
//	m["z"] = 1
//	m["a"] = 2
//	for k := range m {
//	    fmt.Println(k) // Order unpredictable
//	}
//
//	// OrderedMap - predictable insertion order
//	om := omap.New[string, int]()
//	om.Set("z", 1)
//	om.Set("a", 2)
//	for _, k := range om.Keys() {
//	    fmt.Println(k) // Always: z, a
//	}
//
// # When to Use
//
// Use OrderedMap when:
//   - Iteration order matters
//   - Implementing LRU/MRU caches
//   - Maintaining chronological data
//   - Need predictable JSON serialization order
//   - Building configuration or settings managers
//   - Preserving user input order
//
// Use regular map when:
//   - Order doesn't matter
//   - Maximum performance is critical
//   - Memory usage is constrained
//
// # Thread Safety
//
// OrderedMap is not thread-safe. For concurrent access, use external synchronization
// or consider using a concurrent ordered map implementation.
//
//	var mu sync.RWMutex
//	m := omap.New[string, int]()
//
//	mu.Lock()
//	m.Set(key, value)
//	mu.Unlock()
//
//	mu.RLock()
//	val, ok := m.Get(key)
//	mu.RUnlock()
package omap
