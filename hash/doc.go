// Package hash provides efficient, seed-aware hashing utilities for Go types.
//
// It is designed for use in concurrent maps and caches where per-instance
// seeds are required to mitigate hash-flooding attacks.
//
// The package includes optimized hashing for primitive types and structs.
// Struct hashing performs a one-time analysis to compute field offsets,
// enabling fast, allocation-free hashing in performance-critical paths.
package hash
