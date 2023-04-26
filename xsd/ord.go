package xsd

import (
	"reflect"

	"github.com/fogfish/curie"
)

func Compare(a, b Value) int {
	switch av := a.(type) {
	case AnyURI:
		if bv, ok := b.(AnyURI); ok {
			return compare(av.Value, bv.Value)
		} else {
			return compare(reflect.Kind(1000), typeOf(b))
		}
	case String:
		if bv, ok := b.(String); ok {
			return compare(av.Value, bv.Value)
		} else {
			return compare(reflect.String, typeOf(b))
		}
	}
	return 0
}

func typeOf(x any) reflect.Kind {
	switch x.(type) {
	case string:
		return reflect.String
	case bool:
		return reflect.Bool
	case int:
		return reflect.Int
	case int8:
		return reflect.Int8
	case int16:
		return reflect.Int16
	case int32:
		return reflect.Int32
	case int64:
		return reflect.Int64
	case uint:
		return reflect.Uint
	case uint8:
		return reflect.Uint8
	case uint16:
		return reflect.Uint16
	case uint32:
		return reflect.Uint32
	case uint64:
		return reflect.Uint64
	case float32:
		return reflect.Float32
	case float64:
		return reflect.Float64
	case curie.IRI:
		return reflect.Kind(1000)
	case []byte:
		return reflect.Kind(1001)
	default:
		return reflect.Invalid
	}
}

// func value[T comparable](x Value) T {
// 	switch av := x.(type) {
// 	case AnyURI:
// 	case String:
// 		return av.Value
// 	}
// }

// func compareIt[T Value](t reflect.Kind, a T, b any) int {
// 	if bv, ok := b.(T); ok {
// 		return compare(a, bv)
// 	} else {
// 		return compare(t, typeOf(b))
// 	}
// }

func compare[T interface {
	~string |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}](a, b T) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// type Comparable interface {
// 	~string |
// 		~int | ~int8 | ~int16 | ~int32 | ~int64 |
// 		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
// 		~float32 | ~float64
// }
