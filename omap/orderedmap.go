package omap

// OrderedMap is a map that maintains insertion order of keys.
// It combines a hash map for O(1) lookups with a doubly-linked list for order preservation.
type OrderedMap[K comparable, V any] struct {
	items map[K]*entry[K, V]
	head  *entry[K, V]
	tail  *entry[K, V]
	len   int
}

// entry represents a single key-value pair in the ordered map's linked list.
type entry[K comparable, V any] struct {
	key   K
	value V
	prev  *entry[K, V]
	next  *entry[K, V]
}

// Item represents a key-value pair from the ordered map.
type Item[K comparable, V any] struct {
	Key   K
	Value V
}

// New creates and returns a new empty OrderedMap.
func New[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{
		items: make(map[K]*entry[K, V]),
	}
}

// Set inserts or updates a key-value pair.
// If the key already exists, its value is updated and the key is moved to the end.
func (m *OrderedMap[K, V]) Set(key K, value V) {
	if e, exists := m.items[key]; exists {
		e.value = value
		m.moveToBack(e)
		return
	}

	e := &entry[K, V]{
		key:   key,
		value: value,
	}
	m.items[key] = e
	m.addToBack(e)
	m.len++
}

// Get retrieves the value for a key.
// Returns the value and true if found, zero value and false otherwise.
func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	if e, exists := m.items[key]; exists {
		return e.value, true
	}
	var zero V
	return zero, false
}

// Delete removes a key-value pair from the map.
// Returns true if the key was present and removed, false otherwise.
func (m *OrderedMap[K, V]) Delete(key K) bool {
	e, exists := m.items[key]
	if !exists {
		return false
	}

	delete(m.items, key)
	m.removeEntry(e)
	m.len--
	return true
}

// Has checks if a key exists in the map.
func (m *OrderedMap[K, V]) Has(key K) bool {
	_, exists := m.items[key]
	return exists
}

// Len returns the number of key-value pairs in the map.
func (m *OrderedMap[K, V]) Len() int {
	return m.len
}

// Clear removes all key-value pairs from the map.
func (m *OrderedMap[K, V]) Clear() {
	m.items = make(map[K]*entry[K, V])
	m.head = nil
	m.tail = nil
	m.len = 0
}

// Keys returns a slice of all keys in insertion order.
func (m *OrderedMap[K, V]) Keys() []K {
	keys := make([]K, 0, m.len)
	for e := m.head; e != nil; e = e.next {
		keys = append(keys, e.key)
	}
	return keys
}

// Values returns a slice of all values in insertion order.
func (m *OrderedMap[K, V]) Values() []V {
	values := make([]V, 0, m.len)
	for e := m.head; e != nil; e = e.next {
		values = append(values, e.value)
	}
	return values
}

// Items returns a slice of all key-value pairs in insertion order.
func (m *OrderedMap[K, V]) Items() []Item[K, V] {
	items := make([]Item[K, V], 0, m.len)
	for e := m.head; e != nil; e = e.next {
		items = append(items, Item[K, V]{Key: e.key, Value: e.value})
	}
	return items
}

// Range iterates over all key-value pairs in insertion order.
// If the function returns false, iteration stops.
func (m *OrderedMap[K, V]) Range(fn func(K, V) bool) {
	for e := m.head; e != nil; e = e.next {
		if !fn(e.key, e.value) {
			return
		}
	}
}

// First returns the first inserted key-value pair.
// Returns the key, value, and true if the map is not empty, zero values and false otherwise.
func (m *OrderedMap[K, V]) First() (K, V, bool) {
	if m.head == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	return m.head.key, m.head.value, true
}

// Last returns the most recently inserted key-value pair.
// Returns the key, value, and true if the map is not empty, zero values and false otherwise.
func (m *OrderedMap[K, V]) Last() (K, V, bool) {
	if m.tail == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	return m.tail.key, m.tail.value, true
}

// PopFirst removes and returns the first inserted key-value pair.
// Returns the key, value, and true if the map was not empty, zero values and false otherwise.
func (m *OrderedMap[K, V]) PopFirst() (K, V, bool) {
	if m.head == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	e := m.head
	delete(m.items, e.key)
	m.removeEntry(e)
	m.len--
	return e.key, e.value, true
}

// PopLast removes and returns the most recently inserted key-value pair.
// Returns the key, value, and true if the map was not empty, zero values and false otherwise.
func (m *OrderedMap[K, V]) PopLast() (K, V, bool) {
	if m.tail == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	e := m.tail
	delete(m.items, e.key)
	m.removeEntry(e)
	m.len--
	return e.key, e.value, true
}

// Clone creates a deep copy of the OrderedMap with independent internal structures.
// Modifications to the clone will not affect the original map and vice versa.
// The clone preserves the insertion order of the original map.
func (m *OrderedMap[K, V]) Clone() OrderedMap[K, V] {
	clone := New[K, V]()
	for e := m.head; e != nil; e = e.next {
		clone.Set(e.key, e.value)
	}
	return clone
}

// addToBack appends an entry to the end of the linked list.
func (m *OrderedMap[K, V]) addToBack(e *entry[K, V]) {
	if m.tail == nil {
		m.head = e
		m.tail = e
		return
	}

	e.prev = m.tail
	m.tail.next = e
	m.tail = e
}

// moveToBack moves an existing entry to the end of the linked list.
// Used when updating an existing key to maintain most-recently-updated ordering.
func (m *OrderedMap[K, V]) moveToBack(e *entry[K, V]) {
	if e == m.tail {
		return
	}

	m.removeEntry(e)
	m.addToBack(e)
}

// removeEntry removes an entry from the linked list without deleting from the map.
// Updates the prev/next pointers of neighboring entries to maintain list integrity.
func (m *OrderedMap[K, V]) removeEntry(e *entry[K, V]) {
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		m.head = e.next
	}

	if e.next != nil {
		e.next.prev = e.prev
	} else {
		m.tail = e.prev
	}

	e.prev = nil
	e.next = nil
}
