package hash

import (
	"math"
	"strconv"
)

func KeyToBytes[K comparable](key K) []byte {
	switch v := any(key).(type) {
	case string:
		return []byte(v)
	case int:
		return IntToBytes(v)
	case int8:
		return IntToBytes(int(v))
	case int16:
		return IntToBytes(int(v))
	case int32:
		return IntToBytes(int(v))
	case int64:
		return Int64ToBytes(v)
	case uint:
		return UintToBytes(v)
	case uint8:
		return UintToBytes(uint(v))
	case uint16:
		return UintToBytes(uint(v))
	case uint32:
		return UintToBytes(uint(v))
	case uint64:
		return Uint64ToBytes(v)
	case uintptr:
		return Uint64ToBytes(uint64(v))
	case bool:
		if v {
			return []byte{1}
		}
		return []byte{0}
	case float32:
		return Float32ToBytes(v)
	case float64:
		return Float64ToBytes(v)
	case complex64:
		return Complex64ToBytes(v)
	case complex128:
		return Complex128ToBytes(v)
	default:
		return []byte(stringerToString(v))
	}
}

func IntToBytes(n int) []byte {
	return []byte(strconv.Itoa(n))
}

func Int64ToBytes(n int64) []byte {
	return []byte(strconv.FormatInt(n, 10))
}

func UintToBytes(n uint) []byte {
	return []byte(strconv.FormatUint(uint64(n), 10))
}

func Uint64ToBytes(n uint64) []byte {
	return []byte(strconv.FormatUint(n, 10))
}

func Float32ToBytes(f float32) []byte {
	return []byte(strconv.FormatFloat(float64(f), 'f', -1, 32))
}

func Float64ToBytes(f float64) []byte {
	return []byte(strconv.FormatFloat(f, 'f', -1, 64))
}

func Complex64ToBytes(c complex64) []byte {
	realBytes := Float32ToBytes(real(c))
	imagBytes := Float32ToBytes(imag(c))
	return append(realBytes, imagBytes...)
}

func Complex128ToBytes(c complex128) []byte {
	realBytes := Float64ToBytes(real(c))
	imagBytes := Float64ToBytes(imag(c))
	return append(realBytes, imagBytes...)
}

func Float32Bits(f float32) uint32 {
	return math.Float32bits(f)
}

func Float64Bits(f float64) uint64 {
	return math.Float64bits(f)
}

func stringerToString(v any) string {
	type stringer interface {
		String() string
	}
	if s, ok := v.(stringer); ok {
		return s.String()
	}
	return "x"
}
