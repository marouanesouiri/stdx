package cmap

import (
	"fmt"
	"hash/maphash"
	"sync"

	"github.com/marouanesouiri/stdx/hash"
	"github.com/marouanesouiri/stdx/optional"
)

// SHARD_COUNT is the default shard count.
const SHARD_COUNT = 32

// ConcurrentMap is a thread-safe map with high performance through sharding.
// It splits the keyspace across multiple shards, each with its own lock,
// reducing lock contention in concurrent scenarios.
type ConcurrentMap[K comparable, V any] struct {
	shards    []*shard[K, V]
	shardMask uint32
	hashFunc  hash.Hasher[K]
	seed      maphash.Seed
}

// shard represents a single map shard with its own lock.
type shard[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
}

// Item represents a key-value pair from the map.
type Item[K comparable, V any] struct {
	Key   K
	Value V
}

// Option defines a functional option for ConcurrentMap configuration.
type Option[K comparable, V any] func(ConcurrentMap[K, V]) ConcurrentMap[K, V]

// WithHash sets a custom hash function for key sharding.
// The hash function should be fast and provide a good distribution.
func WithHash[K comparable, V any](f hash.Hasher[K]) Option[K, V] {
	return func(m ConcurrentMap[K, V]) ConcurrentMap[K, V] {
		m.hashFunc = f
		return m
	}
}

// WithSeed sets a specific seed for the hash function.
func WithSeed[K comparable, V any](seed maphash.Seed) Option[K, V] {
	return func(m ConcurrentMap[K, V]) ConcurrentMap[K, V] {
		m.seed = seed
		return m
	}
}

// New creates a new ConcurrentMap with default shard count (SHARD_COUNT).
// The shard count is optimized for typical concurrent workloads.
func New[K comparable, V any](opts ...Option[K, V]) ConcurrentMap[K, V] {
	return WithShards(SHARD_COUNT, opts...)
}

// WithShards creates a new ConcurrentMap with the specified number of shards.
// shardCount must be a power of 2 for optimal performance.
// If not a power of 2, it will be rounded up to the next power of 2.
func WithShards[K comparable, V any](shardCount int, opts ...Option[K, V]) ConcurrentMap[K, V] {
	if shardCount <= 0 {
		shardCount = SHARD_COUNT
	}

	shardCount = nextPowerOf2(shardCount)

	shards := make([]*shard[K, V], shardCount)
	for i := range shardCount {
		shards[i] = &shard[K, V]{
			items: make(map[K]V),
		}
	}

	m := ConcurrentMap[K, V]{
		shards:    shards,
		shardMask: uint32(shardCount - 1),
		hashFunc:  hash.GetHashFunc[K](),
		seed:      maphash.MakeSeed(),
	}

	for _, opt := range opts {
		m = opt(m)
	}

	return m
}

// nextPowerOf2 returns the next power of 2 greater than or equal to n.
func nextPowerOf2(n int) int {
	if n <= 0 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}

// getShard returns the shard for the given key.
func (m *ConcurrentMap[K, V]) getShard(key K) *shard[K, V] {
	hashVal := m.hashFunc(m.seed, key)
	index := hashVal & m.shardMask
	return m.shards[index]
}

// Set stores a key-value pair in the map.
func (m *ConcurrentMap[K, V]) Set(key K, value V) {
	shard := m.getShard(key)
	shard.mu.Lock()
	shard.items[key] = value
	shard.mu.Unlock()
}

// Get retrieves a value from the map.
// Returns an Option containing the value if the key exists, otherwise returns None.
func (m *ConcurrentMap[K, V]) Get(key K) optional.Option[V] {
	shard := m.getShard(key)
	shard.mu.RLock()
	val, ok := shard.items[key]
	shard.mu.RUnlock()
	return optional.FromPair(val, ok)
}

// Delete removes a key from the map.
func (m *ConcurrentMap[K, V]) Delete(key K) {
	shard := m.getShard(key)
	shard.mu.Lock()
	delete(shard.items, key)
	shard.mu.Unlock()
}

// Has checks if a key exists in the map.
func (m *ConcurrentMap[K, V]) Has(key K) bool {
	return m.Get(key).IsPresent()
}

// GetOrSet atomically gets a value or sets it if absent.
// Returns the value and true if it existed, or the newly set value and false.
func (m *ConcurrentMap[K, V]) GetOrSet(key K, value V) (V, bool) {
	shard := m.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if existingVal, ok := shard.items[key]; ok {
		return existingVal, true
	}
	shard.items[key] = value
	return value, false
}

// SetIfAbsent sets the value only if the key doesn't exist.
// Returns true if the value was set, false if the key already existed.
func (m *ConcurrentMap[K, V]) SetIfAbsent(key K, value V) bool {
	shard := m.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, ok := shard.items[key]; ok {
		return false
	}
	shard.items[key] = value
	return true
}

// Remove atomically removes and returns a value.
// Returns an Option containing the value if it existed, otherwise returns None.
func (m *ConcurrentMap[K, V]) Remove(key K) optional.Option[V] {
	shard := m.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	val, ok := shard.items[key]
	if ok {
		delete(shard.items, key)
	}
	return optional.FromPair(val, ok)
}

// Compute atomically computes a new value for a key.
// The function receives the current value as an Option.
// The returned value is stored in the map.
func (m *ConcurrentMap[K, V]) Compute(key K, fn func(oldValue optional.Option[V]) V) V {
	shard := m.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	oldValue, exists := shard.items[key]
	newValue := fn(optional.FromPair(oldValue, exists))
	shard.items[key] = newValue
	return newValue
}

// Len returns the total number of items in the map.
func (m *ConcurrentMap[K, V]) Len() int {
	count := 0
	for _, shard := range m.shards {
		shard.mu.RLock()
		count += len(shard.items)
		shard.mu.RUnlock()
	}
	return count
}

// Clear removes all items from the map.
func (m *ConcurrentMap[K, V]) Clear() {
	for _, shard := range m.shards {
		shard.mu.Lock()
		shard.items = make(map[K]V)
		shard.mu.Unlock()
	}
}

// Range calls the function for each key-value pair in the map.
// If the function returns false, iteration stops.
// Note: The function is called while holding a read lock on each shard.
func (m *ConcurrentMap[K, V]) Range(fn func(key K, value V) bool) {
	for _, shard := range m.shards {
		shard.mu.RLock()
		for k, v := range shard.items {
			if !fn(k, v) {
				shard.mu.RUnlock()
				return
			}
		}
		shard.mu.RUnlock()
	}
}

// Keys returns a slice of all keys in the map.
// This creates a snapshot at the time of the call.
func (m *ConcurrentMap[K, V]) Keys() []K {
	keys := make([]K, 0, m.Len())
	m.Range(func(key K, _ V) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

// Values returns a slice of all values in the map.
// This creates a snapshot at the time of the call.
func (m *ConcurrentMap[K, V]) Values() []V {
	values := make([]V, 0, m.Len())
	m.Range(func(_ K, value V) bool {
		values = append(values, value)
		return true
	})
	return values
}

// Items returns a slice of all key-value pairs in the map.
// This creates a snapshot at the time of the call.
func (m *ConcurrentMap[K, V]) Items() []Item[K, V] {
	items := make([]Item[K, V], 0, m.Len())
	m.Range(func(key K, value V) bool {
		items = append(items, Item[K, V]{Key: key, Value: value})
		return true
	})
	return items
}

// Clone creates a deep copy of the ConcurrentMap with independent shards.
// Modifications to the clone will not affect the original map and vice versa.
// This operation locks all shards temporarily to ensure a consistent snapshot.
func (m *ConcurrentMap[K, V]) Clone() ConcurrentMap[K, V] {
	clone := WithShards(len(m.shards), WithHash[K, V](m.hashFunc), WithSeed[K, V](m.seed))
	m.Range(func(key K, value V) bool {
		clone.Set(key, value)
		return true
	})
	return clone
}

// String returns a string representation of this cmap.
func (m *ConcurrentMap[K, V]) String() string {
	return fmt.Sprintf("ConcurrentMap{len=%d, shards=%d}", m.Len(), len(m.shards))
}
