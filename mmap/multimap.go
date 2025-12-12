package mmap

// Multimap is a map that allows multiple values per key.
// It prevents duplicate values for the same key.
type Multimap[K comparable, V comparable] struct {
	items map[K]map[V]struct{}
	size  int
}

// Entry represents a single key-value pair from the multimap.
type Entry[K comparable, V comparable] struct {
	Key   K
	Value V
}

// New creates and returns a new empty Multimap.
func New[K comparable, V comparable]() Multimap[K, V] {
	return Multimap[K, V]{
		items: make(map[K]map[V]struct{}),
	}
}

// Put adds a value to the set of values for a key.
// Returns true if the value was added, false if it already existed.
func (m *Multimap[K, V]) Put(key K, value V) bool {
	if m.items[key] == nil {
		m.items[key] = make(map[V]struct{})
	}

	if _, exists := m.items[key][value]; exists {
		return false
	}

	m.items[key][value] = struct{}{}
	m.size++
	return true
}

// PutAll adds multiple values for a key.
// Returns the count of values that were actually added (excludes duplicates).
func (m *Multimap[K, V]) PutAll(key K, values ...V) int {
	count := 0
	for _, v := range values {
		if m.Put(key, v) {
			count++
		}
	}
	return count
}

// Get returns all values associated with a key.
// Returns an empty slice if the key doesn't exist.
func (m *Multimap[K, V]) Get(key K) []V {
	set, exists := m.items[key]
	if !exists {
		return []V{}
	}

	values := make([]V, 0, len(set))
	for v := range set {
		values = append(values, v)
	}
	return values
}

// GetSet returns the set of values for a key as a map.
// Returns an empty map if the key doesn't exist.
func (m *Multimap[K, V]) GetSet(key K) map[V]struct{} {
	if set, exists := m.items[key]; exists {
		result := make(map[V]struct{}, len(set))
		for v := range set {
			result[v] = struct{}{}
		}
		return result
	}
	return make(map[V]struct{})
}

// Delete removes a specific value for a key.
// Returns true if the value was present and removed, false otherwise.
func (m *Multimap[K, V]) Delete(key K, value V) bool {
	set, exists := m.items[key]
	if !exists {
		return false
	}

	if _, hasValue := set[value]; !hasValue {
		return false
	}

	delete(set, value)
	m.size--

	if len(set) == 0 {
		delete(m.items, key)
	}

	return true
}

// DeleteAll removes all values for a key.
// Returns true if the key existed, false otherwise.
func (m *Multimap[K, V]) DeleteAll(key K) bool {
	set, exists := m.items[key]
	if !exists {
		return false
	}

	m.size -= len(set)
	delete(m.items, key)
	return true
}

// Contains checks if a specific key-value pair exists.
func (m *Multimap[K, V]) Contains(key K, value V) bool {
	set, exists := m.items[key]
	if !exists {
		return false
	}
	_, hasValue := set[value]
	return hasValue
}

// ContainsKey checks if a key exists in the multimap.
func (m *Multimap[K, V]) ContainsKey(key K) bool {
	_, exists := m.items[key]
	return exists
}

// Size returns the total number of key-value pairs.
func (m *Multimap[K, V]) Size() int {
	return m.size
}

// KeySize returns the number of values for a specific key.
func (m *Multimap[K, V]) KeySize(key K) int {
	if set, exists := m.items[key]; exists {
		return len(set)
	}
	return 0
}

// Len returns the number of unique keys.
func (m *Multimap[K, V]) Len() int {
	return len(m.items)
}

// Clear removes all key-value pairs from the multimap.
func (m *Multimap[K, V]) Clear() {
	m.items = make(map[K]map[V]struct{})
	m.size = 0
}

// Keys returns a slice of all unique keys.
func (m *Multimap[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.items))
	for k := range m.items {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of all values across all keys.
func (m *Multimap[K, V]) Values() []V {
	values := make([]V, 0, m.size)
	for _, set := range m.items {
		for v := range set {
			values = append(values, v)
		}
	}
	return values
}

// Entries returns all key-value pairs as a slice.
func (m *Multimap[K, V]) Entries() []Entry[K, V] {
	entries := make([]Entry[K, V], 0, m.size)
	for k, set := range m.items {
		for v := range set {
			entries = append(entries, Entry[K, V]{Key: k, Value: v})
		}
	}
	return entries
}

// Range iterates over all key-value pairs.
// If the function returns false, iteration stops.
func (m *Multimap[K, V]) Range(fn func(K, V) bool) {
	for k, set := range m.items {
		for v := range set {
			if !fn(k, v) {
				return
			}
		}
	}
}

// ForEachKey iterates over keys with their associated values.
// If the function returns false, iteration stops.
func (m *Multimap[K, V]) ForEachKey(fn func(K, []V) bool) {
	for k := range m.items {
		values := m.Get(k)
		if !fn(k, values) {
			return
		}
	}
}
