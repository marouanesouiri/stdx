package cmap

import (
	"sync"
	"testing"

	"github.com/marouanesouiri/stdx/optional"
)

// TestConcurrentMapBasic tests basic operations
func TestConcurrentMapBasic(t *testing.T) {
	m := New[string, int]()

	// Test Set and Get
	m.Set("key1", 100)
	if opt := m.Get("key1"); !opt.IsPresent() || opt.MustGet() != 100 {
		t.Errorf("Expected 100, got %v", opt)
	}

	// Test Has
	if !m.Has("key1") {
		t.Error("Expected key1 to exist")
	}

	// Test Delete
	m.Delete("key1")
	if m.Has("key1") {
		t.Error("Expected key1 to be deleted")
	}
}

// TestConcurrentMapConcurrency tests thread safety
func TestConcurrentMapConcurrency(t *testing.T) {
	m := New[int, int]()
	const goroutines = 100
	const operations = 1000

	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := id*operations + j
				m.Set(key, key)
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := id*operations + j
				m.Get(key)
			}
		}(i)
	}

	wg.Wait()

	// Verify total count
	expectedCount := goroutines * operations
	if m.Len() != expectedCount {
		t.Errorf("Expected %d items, got %d", expectedCount, m.Len())
	}
}

// TestConcurrentMapAtomic tests atomic operations
func TestConcurrentMapAtomic(t *testing.T) {
	m := New[string, int]()

	// Test GetOrSet
	val, existed := m.GetOrSet("counter", 1)
	if existed || val != 1 {
		t.Error("Expected new value to be set")
	}

	val, existed = m.GetOrSet("counter", 10)
	if !existed || val != 1 {
		t.Error("Expected existing value to be returned")
	}

	// Test SetIfAbsent
	if !m.SetIfAbsent("new", 42) {
		t.Error("Expected SetIfAbsent to succeed")
	}
	if m.SetIfAbsent("new", 100) {
		t.Error("Expected SetIfAbsent to fail on existing key")
	}

	// Test Remove
	m.Set("temp", 99)
	if opt := m.Remove("temp"); !opt.IsPresent() || opt.MustGet() != 99 {
		t.Error("Expected Remove to return value")
	}
	if opt := m.Remove("temp"); opt.IsPresent() {
		t.Error("Expected Remove to fail on missing key")
	}

	// Test Compute
	m.Set("compute", 5)
	newVal := m.Compute("compute", func(old optional.Option[int]) int {
		if !old.IsPresent() {
			return 1
		}
		return old.MustGet() * 2
	})
	if newVal != 10 {
		t.Errorf("Expected 10, got %d", newVal)
	}
}

// TestConcurrentMapIteration tests Range, Keys, Values
func TestConcurrentMapIteration(t *testing.T) {
	m := New[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	// Test Range
	count := 0
	m.Range(func(key string, value int) bool {
		count++
		return true
	})
	if count != 3 {
		t.Errorf("Expected 3 items in Range, got %d", count)
	}

	// Test Keys
	keys := m.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Test Values
	values := m.Values()
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	// Test Items
	items := m.Items()
	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	// Test Clear
	m.Clear()
	if m.Len() != 0 {
		t.Errorf("Expected 0 items after Clear, got %d", m.Len())
	}
}

// BenchmarkConcurrentMapSet benchmarks Set operations
func BenchmarkConcurrentMapSet(b *testing.B) {
	m := New[int, int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Set(i, i)
			i++
		}
	})
}

// BenchmarkConcurrentMapGet benchmarks Get operations
func BenchmarkConcurrentMapGet(b *testing.B) {
	m := New[int, int]()
	for i := 0; i < 10000; i++ {
		m.Set(i, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Get(i % 10000)
			i++
		}
	})
}
