package hash

import (
	"encoding/binary"
	"hash/maphash"
	"math"
	"reflect"
	"unsafe"
)

// Hashable is an interface for types that can hash themselves.
// This allows for high-performance, domain-specific hashing.
type Hashable interface {
	Hash(seed maphash.Seed) uint32
}

// Hasher is a function type that takes a seed and a value, and returns its hash.
type Hasher[T any] func(maphash.Seed, T) uint32

// StringHasher returns a hash for the given string using the provided seed.
func StringHasher(seed maphash.Seed, s string) uint32 {
	return uint32(maphash.String(seed, s))
}

// IntHasher returns a hash for the given int using the provided seed.
func IntHasher(seed maphash.Seed, n int) uint32 {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(n))
	return uint32(maphash.Bytes(seed, b[:]))
}

// Int8Hasher returns a hash for the given int8 using the provided seed.
func Int8Hasher(seed maphash.Seed, n int8) uint32 {
	return uint32(maphash.Bytes(seed, []byte{byte(n)}))
}

// Int16Hasher returns a hash for the given int16 using the provided seed.
func Int16Hasher(seed maphash.Seed, n int16) uint32 {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], uint16(n))
	return uint32(maphash.Bytes(seed, b[:]))
}

// Int32Hasher returns a hash for the given int32 using the provided seed.
func Int32Hasher(seed maphash.Seed, n int32) uint32 {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(n))
	return uint32(maphash.Bytes(seed, b[:]))
}

// Int64Hasher returns a hash for the given int64 using the provided seed.
func Int64Hasher(seed maphash.Seed, n int64) uint32 {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(n))
	return uint32(maphash.Bytes(seed, b[:]))
}

// UintHasher returns a hash for the given uint using the provided seed.
func UintHasher(seed maphash.Seed, n uint) uint32 {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(n))
	return uint32(maphash.Bytes(seed, b[:]))
}

// Uint8Hasher returns a hash for the given uint8 using the provided seed.
func Uint8Hasher(seed maphash.Seed, n uint8) uint32 {
	return uint32(maphash.Bytes(seed, []byte{n}))
}

// Uint16Hasher returns a hash for the given uint16 using the provided seed.
func Uint16Hasher(seed maphash.Seed, n uint16) uint32 {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], n)
	return uint32(maphash.Bytes(seed, b[:]))
}

// Uint32Hasher returns a hash for the given uint32 using the provided seed.
func Uint32Hasher(seed maphash.Seed, n uint32) uint32 {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], n)
	return uint32(maphash.Bytes(seed, b[:]))
}

// Uint64Hasher returns a hash for the given uint64 using the provided seed.
func Uint64Hasher(seed maphash.Seed, n uint64) uint32 {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], n)
	return uint32(maphash.Bytes(seed, b[:]))
}

// UintptrHasher returns a hash for the given uintptr using the provided seed.
func UintptrHasher(seed maphash.Seed, n uintptr) uint32 {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(n))
	return uint32(maphash.Bytes(seed, b[:]))
}

// Float32Hasher returns a hash for the given float32 using the provided seed.
func Float32Hasher(seed maphash.Seed, f float32) uint32 {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], math.Float32bits(f))
	return uint32(maphash.Bytes(seed, b[:]))
}

// Float64Hasher returns a hash for the given float64 using the provided seed.
func Float64Hasher(seed maphash.Seed, f float64) uint32 {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], math.Float64bits(f))
	return uint32(maphash.Bytes(seed, b[:]))
}

// BoolHasher returns a hash for the given bool using the provided seed.
func BoolHasher(seed maphash.Seed, v bool) uint32 {
	val := byte(0)
	if v {
		val = 1
	}
	return uint32(maphash.Bytes(seed, []byte{val}))
}

// GetHashFunc returns a hash function for the given comparable type K.
// The returned function takes a seed and a value.
func GetHashFunc[K comparable]() Hasher[K] {
	var k K
	switch any(k).(type) {
	case string:
		return func(seed maphash.Seed, key K) uint32 {
			return StringHasher(seed, any(key).(string))
		}
	case int:
		return func(seed maphash.Seed, key K) uint32 {
			return IntHasher(seed, any(key).(int))
		}
	case int8:
		return func(seed maphash.Seed, key K) uint32 {
			return Int8Hasher(seed, any(key).(int8))
		}
	case int16:
		return func(seed maphash.Seed, key K) uint32 {
			return Int16Hasher(seed, any(key).(int16))
		}
	case int32:
		return func(seed maphash.Seed, key K) uint32 {
			return Int32Hasher(seed, any(key).(int32))
		}
	case int64:
		return func(seed maphash.Seed, key K) uint32 {
			return Int64Hasher(seed, any(key).(int64))
		}
	case uint:
		return func(seed maphash.Seed, key K) uint32 {
			return UintHasher(seed, any(key).(uint))
		}
	case uint8:
		return func(seed maphash.Seed, key K) uint32 {
			return Uint8Hasher(seed, any(key).(uint8))
		}
	case uint16:
		return func(seed maphash.Seed, key K) uint32 {
			return Uint16Hasher(seed, any(key).(uint16))
		}
	case uint32:
		return func(seed maphash.Seed, key K) uint32 {
			return Uint32Hasher(seed, any(key).(uint32))
		}
	case uint64:
		return func(seed maphash.Seed, key K) uint32 {
			return Uint64Hasher(seed, any(key).(uint64))
		}
	case uintptr:
		return func(seed maphash.Seed, key K) uint32 {
			return UintptrHasher(seed, any(key).(uintptr))
		}
	case float32:
		return func(seed maphash.Seed, key K) uint32 {
			return Float32Hasher(seed, any(key).(float32))
		}
	case float64:
		return func(seed maphash.Seed, key K) uint32 {
			return Float64Hasher(seed, any(key).(float64))
		}
	case bool:
		return func(seed maphash.Seed, key K) uint32 {
			return BoolHasher(seed, any(key).(bool))
		}
	default:
		t := reflect.TypeOf(k)
		if t.Kind() == reflect.Struct {
			return CreateStructHasher[K](t)
		}

		return func(seed maphash.Seed, key K) uint32 {
			if h, ok := any(key).(Hashable); ok {
				return h.Hash(seed)
			}
			var h maphash.Hash
			h.SetSeed(seed)
			switch v := any(key).(type) {
			case interface{ String() string }:
				h.WriteString(v.String())
			default:
			}
			return uint32(h.Sum64())
		}
	}
}

type fieldInfo struct {
	offset uintptr
	kind   reflect.Kind
}

// CreateStructHasher returns a Hasher func for the struct K.
// The hasher is a simple hash func that uses Fibonacci Hashing.
func CreateStructHasher[K comparable](t reflect.Type) Hasher[K] {
	fields := flattenStruct(t, 0)

	return func(seed maphash.Seed, key K) uint32 {
		var h uint32
		p := unsafe.Pointer(&key)
		for _, f := range fields {
			fieldPtr := unsafe.Pointer(uintptr(p) + f.offset)
			var fHash uint32
			switch f.kind {
			case reflect.String:
				fHash = StringHasher(seed, *(*string)(fieldPtr))
			case reflect.Int:
				fHash = IntHasher(seed, *(*int)(fieldPtr))
			case reflect.Int8:
				fHash = Int8Hasher(seed, *(*int8)(fieldPtr))
			case reflect.Int16:
				fHash = Int16Hasher(seed, *(*int16)(fieldPtr))
			case reflect.Int32:
				fHash = Int32Hasher(seed, *(*int32)(fieldPtr))
			case reflect.Int64:
				fHash = Int64Hasher(seed, *(*int64)(fieldPtr))
			case reflect.Uint:
				fHash = UintHasher(seed, *(*uint)(fieldPtr))
			case reflect.Uint8:
				fHash = Uint8Hasher(seed, *(*uint8)(fieldPtr))
			case reflect.Uint16:
				fHash = Uint16Hasher(seed, *(*uint16)(fieldPtr))
			case reflect.Uint32:
				fHash = Uint32Hasher(seed, *(*uint32)(fieldPtr))
			case reflect.Uint64:
				fHash = Uint64Hasher(seed, *(*uint64)(fieldPtr))
			case reflect.Uintptr:
				fHash = UintptrHasher(seed, *(*uintptr)(fieldPtr))
			case reflect.Float32:
				fHash = Float32Hasher(seed, *(*float32)(fieldPtr))
			case reflect.Float64:
				fHash = Float64Hasher(seed, *(*float64)(fieldPtr))
			case reflect.Bool:
				fHash = BoolHasher(seed, *(*bool)(fieldPtr))
			case reflect.Pointer, reflect.UnsafePointer:
				var b [8]byte
				binary.LittleEndian.PutUint64(b[:], uint64(uintptr(*(*unsafe.Pointer)(fieldPtr))))
				fHash = uint32(maphash.Bytes(seed, b[:]))
			}
			h ^= fHash + 0x9e3779b9 + (h << 6) + (h >> 2)
		}
		return h
	}
}

func flattenStruct(st reflect.Type, baseOffset uintptr) []fieldInfo {
	var fields []fieldInfo
	for i := range st.NumField() {
		f := st.Field(i)
		if f.Name == "_" {
			continue
		}
		kind := f.Type.Kind()
		switch kind {
		case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64, reflect.Bool, reflect.Pointer, reflect.UnsafePointer:
			fields = append(fields, fieldInfo{offset: baseOffset + f.Offset, kind: kind})
		case reflect.Struct:
			fields = append(fields, flattenStruct(f.Type, baseOffset+f.Offset)...)
		default:
			continue
		}
	}
	return fields
}
