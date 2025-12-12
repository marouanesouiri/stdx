// Package mmap provides a multimap data structure that maps keys to multiple values.
//
// Multimap allows a single key to be associated with multiple values.
// Unlike a regular map where each key maps to a single value, a Multimap can store
// multiple values per key. This is similar to Java's Guava Multimap.
//
// # Basic Usage
//
// Create and use a multimap:
//
//	m := mmap.New[string, string]()
//	m.Put("tags", "go")
//	m.Put("tags", "generics")
//	m.Put("tags", "stdlib")
//
//	values := m.Get("tags") // ["go", "generics", "stdlib"]
//
// # Adding Values
//
//	m := mmap.New[string, int]()
//
//	added := m.Put("numbers", 1)    // true (new value)
//	added = m.Put("numbers", 2)     // true (new value)
//	added = m.Put("numbers", 1)     // false (duplicate)
//
//	count := m.PutAll("more", 3, 4, 5) // Adds 3 values, returns 3
//
// # Retrieving Values
//
//	m := mmap.New[string, string]()
//	m.Put("colors", "red")
//	m.Put("colors", "blue")
//
//	values := m.Get("colors")    // []string{"red", "blue"}
//	set := m.GetSet("colors")    // map[string]struct{}{"red":{}, "blue":{}}
//	empty := m.Get("nonexistent") // []string{} (empty slice)
//
// # Checking Existence
//
//	m := mmap.New[string, int]()
//	m.Put("key", 1)
//	m.Put("key", 2)
//
//	hasKey := m.ContainsKey("key")   // true
//	hasEntry := m.Contains("key", 1) // true
//	hasEntry = m.Contains("key", 3)  // false
//
// # Removing Values
//
//	m := mmap.New[string, int]()
//	m.Put("nums", 1)
//	m.Put("nums", 2)
//	m.Put("nums", 3)
//
//	removed := m.Delete("nums", 2)  // true (removes only value 2)
//	values := m.Get("nums")         // [1, 3]
//
//	removed = m.DeleteAll("nums")   // true (removes all values for key)
//	hasKey := m.ContainsKey("nums") // false
//
// # Size Information
//
//	m := mmap.New[string, int]()
//	m.Put("a", 1)
//	m.Put("a", 2)
//	m.Put("b", 3)
//
//	totalPairs := m.Size()    // 3 (total key-value pairs)
//	numKeys := m.Len()        // 2 (number of unique keys)
//	keySize := m.KeySize("a") // 2 (values for key "a")
//
// # Iteration
//
// Iterate over all key-value pairs:
//
//	m := mmap.New[string, string]()
//	m.Put("fruit", "apple")
//	m.Put("fruit", "banana")
//	m.Put("veggie", "carrot")
//
//	m.Range(func(key, value string) bool {
//	    fmt.Printf("%s: %s\n", key, value)
//	    return true // Continue iteration
//	})
//
// Iterate by key:
//
//	m.ForEachKey(func(key string, values []string) bool {
//	    fmt.Printf("%s has %d values\n", key, len(values))
//	    return true
//	})
//
// # Getting Collections
//
//	m := mmap.New[string, int]()
//	m.Put("a", 1)
//	m.Put("a", 2)
//	m.Put("b", 3)
//
//	keys := m.Keys()       // []string{"a", "b"}
//	values := m.Values()   // []int{1, 2, 3}
//	entries := m.Entries() // []Entry{{Key:"a",Value:1}, ...}
//
// # Use Cases
//
// **Tagging System:**
//
//	tags := mmap.New[string, string]()
//	tags.Put("article-1", "go")
//	tags.Put("article-1", "programming")
//	tags.Put("article-1", "tutorial")
//	tags.Put("article-2", "go")
//	tags.Put("article-2", "advanced")
//
//	articleTags := tags.Get("article-1")
//
// **Reverse Index:**
//
//	index := mmap.New[string, int]()
//	index.Put("golang", 1)
//	index.Put("golang", 5)
//	index.Put("tutorial", 1)
//	index.Put("tutorial", 3)
//	index.Put("tutorial", 5)
//
//	docIDs := index.Get("golang") // Documents containing "golang"
//
// **Graph Adjacency List:**
//
//	graph := mmap.New[string, string]()
//	graph.Put("A", "B")
//	graph.Put("A", "C")
//	graph.Put("B", "C")
//	graph.Put("B", "D")
//
//	neighbors := graph.Get("A") // ["B", "C"]
//
// **Grouping Data:**
//
//	users := mmap.New[string, User]()
//	users.Put("admin", alice)
//	users.Put("admin", bob)
//	users.Put("user", charlie)
//
//	admins := users.Get("admin")
//
// # Performance
//
// **Time Complexity:**
//   - Put: O(1) average
//   - Get: O(k) where k is values per key
//   - Delete: O(1) average
//   - DeleteAll: O(k) where k is values per key
//   - Contains: O(1) average
//   - Size/Len: O(1)
//
// **Space Complexity:**
//   - O(n) where n is total key-value pairs
//   - Additional overhead: One set (map) per unique key
//
// # Comparison with map[K][]V
//
//	// Regular map with slice - allows duplicates, manual management
//	m := make(map[string][]int)
//	m["key"] = append(m["key"], 1)
//	m["key"] = append(m["key"], 1) // Duplicate allowed
//
//	// Multimap - no duplicates, automatic set management
//	mm := mmap.New[string, int]()
//	mm.Put("key", 1)     // true
//	mm.Put("key", 1)     // false (duplicate prevented)
//
// # When to Use
//
// Use Multimap when:
//   - One key needs multiple distinct values
//   - Implementing tags, labels, or categories
//   - Building reverse indexes
//   - Representing graph adjacency lists
//   - Grouping related items
//   - Preventing duplicate values per key
//
// Use map[K][]V when:
//   - Need to preserve value order
//   - Allow duplicate values
//   - Values are accessed as a whole slice
//
// # No Duplicates
//
// Multimap uses sets internally, so duplicate values per key are automatically prevented:
//
//	m := mmap.New[string, int]()
//	m.Put("nums", 1) // true (added)
//	m.Put("nums", 2) // true (added)
//	m.Put("nums", 1) // false (duplicate, not added)
//
//	values := m.Get("nums") // [1, 2] (no duplicate 1)
//
// # Thread Safety
//
// Multimap is not thread-safe. For concurrent access, use external synchronization.
//
//	var mu sync.RWMutex
//	m := mmap.New[string, int]()
//
//	mu.Lock()
//	m.Put(key, value)
//	mu.Unlock()
//
//	mu.RLock()
//	values := m.Get(key)
//	mu.RUnlock()
//
// # Value Type Constraints
//
// Both K and V must be comparable. This is required because values are stored in a set (map).
// If you need non-comparable values, consider using map[K][]V directly.
package mmap
