//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/hexagon
//

package ord

import (
	"reflect"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid"
	"github.com/fogfish/skiplist/ord"
)

const (
	Any = tAny("hexagon.ord.any")
	LID = tLID("hexagon.ord.lid")
	IRI = tIRI("hexagon.ord.iri")
)

// IRI is an instance of ord.Ord[curie.IRI] type class
type tIRI string

func (tIRI) Compare(a, b curie.IRI) int {
	return ord.String.Compare(string(a), string(b))
}

var _ ord.Ord[curie.IRI] = tIRI("")

// LID is an instance of ord.Ord[guid.LID] type class
type tLID string

func (tLID) Compare(a, b guid.LID) int {
	if guid.L.Less(a, b) {
		return -1
	}

	if guid.L.Equal(a, b) {
		return 0
	}

	return 1
}

var _ ord.Ord[guid.LID] = tLID("")

//
type tAny string

func (ord tAny) Compare(a, b any) int {
	switch av := a.(type) {
	case string:
		return compareIt(reflect.String, av, b)
	case bool:
		if bv, ok := b.(bool); ok {
			if av == bv {
				return 0
			}

			if !av && bv {
				return -1
			}

			return 1
		} else {
			return compare(reflect.Bool, typeOf(b))
		}
	case int:
		return compareIt(reflect.Int, av, b)
	case int8:
		return compareIt(reflect.Int8, av, b)
	case int16:
		return compareIt(reflect.Int16, av, b)
	case int32:
		return compareIt(reflect.Int32, av, b)
	case int64:
		return compareIt(reflect.Int64, av, b)
	case uint:
		return compareIt(reflect.Uint, av, b)
	case uint8:
		return compareIt(reflect.Uint8, av, b)
	case uint16:
		return compareIt(reflect.Uint16, av, b)
	case uint32:
		return compareIt(reflect.Uint32, av, b)
	case uint64:
		return compareIt(reflect.Uint64, av, b)
	case float32:
		return compareIt(reflect.Float32, av, b)
	case float64:
		return compareIt(reflect.Float64, av, b)
	case curie.IRI:
		return compareIt(reflect.Kind(1000), av, b)
	case []byte:
		if bv, ok := b.([]byte); ok {
			return compare(string(av), string(bv))
		} else {
			return 1
		}
	default:
		return 1
	}
}

var _ ord.Ord[any] = tAny("")

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

func compareIt[T ord.Comparable](t reflect.Kind, a T, b any) int {
	if bv, ok := b.(T); ok {
		return compare(a, bv)
	} else {
		return compare(t, typeOf(b))
	}
}

func compare[T ord.Comparable](a, b T) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
