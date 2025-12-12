// Package set provides a generic Set data structure for storing unique elements.
//
// Set is a collection that contains no duplicate elements. It models the mathematical
// set abstraction and provides operations for testing membership, computing unions,
// intersections, and differences.
//
// # Basic Usage
//
// Create and use a set:
//
//	s := set.New[string]()
//	s.Add("apple")
//	s.Add("banana")
//	s.Add("apple") // Duplicate, not added
//
//	fmt.Println(s.Size())         // 2
//	fmt.Println(s.Contains("apple")) // true
//
// Create from a slice:
//
//	numbers := []int{1, 2, 3, 2, 1}
//	s := set.FromSlice(numbers)
//	fmt.Println(s.Size()) // 3 (duplicates removed)
//
// # Adding and Removing
//
//	s := set.New[int]()
//
//	added := s.Add(1)        // true (new element)
//	added = s.Add(1)         // false (duplicate)
//
//	count := s.AddAll(2, 3, 4, 3) // Adds 3 new elements
//	fmt.Println(count)       // 3
//
//	removed := s.Remove(2)   // true
//	removed = s.Remove(99)   // false (not present)
//
// # Membership Testing
//
//	s := set.New[string]()
//	s.Add("hello")
//
//	if s.Contains("hello") {
//	    fmt.Println("Found!")
//	}
//
//	if s.IsEmpty() {
//	    fmt.Println("Set is empty")
//	}
//
//	size := s.Size()
//
// # Set Operations
//
// **Union** - All elements from both sets:
//
//	s1 := set.FromSlice([]int{1, 2, 3})
//	s2 := set.FromSlice([]int{3, 4, 5})
//	union := s1.Union(s2)
//	fmt.Println(union.ToSlice()) // [1, 2, 3, 4, 5]
//
// **Intersection** - Elements in both sets:
//
//	intersection := s1.Intersection(s2)
//	fmt.Println(intersection.ToSlice()) // [3]
//
// **Difference** - Elements in first set but not second:
//
//	diff := s1.Difference(s2)
//	fmt.Println(diff.ToSlice()) // [1, 2]
//
// **Symmetric Difference** - Elements in either set but not both:
//
//	symDiff := s1.SymmetricDifference(s2)
//	fmt.Println(symDiff.ToSlice()) // [1, 2, 4, 5]
//
// # Subset and Superset
//
//	s1 := set.FromSlice([]int{1, 2})
//	s2 := set.FromSlice([]int{1, 2, 3, 4})
//
//	if s1.IsSubset(s2) {
//	    fmt.Println("s1 is a subset of s2")
//	}
//
//	if s2.IsSuperset(s1) {
//	    fmt.Println("s2 is a superset of s1")
//	}
//
// # Equality
//
//	s1 := set.FromSlice([]int{1, 2, 3})
//	s2 := set.FromSlice([]int{3, 2, 1})
//
//	if s1.Equal(s2) {
//	    fmt.Println("Sets are equal") // Order doesn't matter
//	}
//
// # Iteration
//
//	s := set.FromSlice([]string{"a", "b", "c"})
//
//	s.Range(func(item string) bool {
//	    fmt.Println(item)
//	    return true // Continue iteration
//	})
//
//	slice := s.ToSlice() // Convert to slice
//
// # Copying Sets
//
//	original := set.FromSlice([]int{1, 2, 3})
//	copy := original.Clone()
//
//	copy.Add(4)
//	fmt.Println(original.Size()) // 3 (unchanged)
//	fmt.Println(copy.Size())     // 4
//
// You can also copy by assignment since Set uses value semantics:
//
//	s1 := set.New[int]()
//	s1.Add(1)
//	s2 := s1 // Both reference the same underlying map
//	s2.Add(2)
//	fmt.Println(s1.Size()) // 2 (same map)
//
// # Use Cases
//
// **Removing Duplicates:**
//
//	data := []int{1, 2, 2, 3, 3, 3, 4}
//	unique := set.FromSlice(data).ToSlice()
//
// **Tracking Seen Items:**
//
//	seen := set.New[string]()
//	for _, item := range items {
//	    if seen.Contains(item) {
//	        fmt.Println("Duplicate:", item)
//	        continue
//	    }
//	    seen.Add(item)
//	    process(item)
//	}
//
// **Set Algebra:**
//
//	admins := set.FromSlice([]string{"alice", "bob"})
//	users := set.FromSlice([]string{"bob", "charlie"})
//
//	allPeople := admins.Union(users)
//	bothRoles := admins.Intersection(users)
//	adminOnly := admins.Difference(users)
//
// **Tag Filtering:**
//
//	required := set.FromSlice([]string{"go", "backend"})
//	articleTags := set.FromSlice([]string{"go", "backend", "tutorial"})
//
//	if required.IsSubset(articleTags) {
//	    fmt.Println("Article matches all required tags")
//	}
//
// # Performance
//
// **Time Complexity:**
//   - Add: O(1) average
//   - Remove: O(1) average
//   - Contains: O(1) average
//   - Size/IsEmpty: O(1)
//   - Union/Intersection/Difference: O(n+m)
//   - IsSubset/IsSuperset: O(n)
//   - Equal: O(n)
//
// **Space Complexity:**
//   - O(n) where n is the number of elements
//   - Uses map[T]struct{} for zero-byte values
//
// # Value Semantics
//
// Set uses value semantics (not pointer receivers). This is safe and efficient because:
//   - The underlying map is a reference type
//   - Copying a Set copies only the map header (~24 bytes), not the data
//   - Multiple Set values can share the same underlying map data
//
// For independent sets, use Clone():
//
//	s1 := set.New[int]()
//	s2 := s1.Clone() // Independent copy
//
// # Comparison with map[T]struct{}
//
//	// Manual set with map
//	m := make(map[string]struct{})
//	m["item"] = struct{}{}
//	_, exists := m["item"]
//
//	// Set package
//	s := set.New[string]()
//	s.Add("item")
//	exists := s.Contains("item")
//
// Set provides:
//   - Cleaner API
//   - Set operations (union, intersection, etc.)
//   - Convenience methods
//   - Clear intent in code
//
// # Thread Safety
//
// Set is not thread-safe. For concurrent access, use external synchronization:
//
//	var mu sync.RWMutex
//	s := set.New[int]()
//
//	mu.Lock()
//	s.Add(1)
//	mu.Unlock()
//
//	mu.RLock()
//	exists := s.Contains(1)
//	mu.RUnlock()
package set
