//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/hexagon
//

package hexagon

import (
	"math/rand"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid"
	"github.com/fogfish/hexagon/internal/ord"
	"github.com/fogfish/skiplist"
)

/*

DataType is a type constrain used by the library.
See https://www.w3.org/TR/xmlschema-2/#datatype

Knowledge statements contain scalar objects -- literals.
Literals are either language-tagged string `rdf:langString` or type-safe values
containing a reference to data-type (e.g. xsd:string).

This interface defines data-types supported by the library. It maps well-known
semantic types to Golang native types and relation to existed schema(s) and
ontologies.

xsd:anyURI ⇒ curie.IRI

The data type represents Internationalized Resource Identifier.
Used to uniquely identify concept, objects, etc.

xsd:string ⇒ string

The string data-type represents character strings in knowledge statements.
The language strings are annotated with corresponding language tag.

xsd:integer ⇒ int

The library uses various int precision data-types to represent decimal values.

xsd:nonNegativeInteger ⇒ uint
xsd:unsignedByte ⇒ uint8
xsd:unsignedShort ⇒ uint16
xsd:unsignedInt ⇒ uint32
xsd:unsignedLong ⇒ uint64

The library uses various uint precision data-types to represent positive decimal values.

xsd:float ⇒ float32
xsd:double ⇒ float64

The value is the IEEE 754 double-precision 64-bit floating point type.

xsd:boolean ⇒ bool

The value is either true or false, representing a logic values

xsd:hexBinary ⇒ []byte
xsd:base64Binary ⇒ []byte

*/
type DataType interface {
	~string |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~bool |
		~[]byte
}

/*

Stream of <s,p,o> triples fetched from the store
*/
type Stream interface {
	Head() (curie.IRI, curie.IRI, any)
	Next() bool
	FMap(f func(curie.IRI, curie.IRI, any) error) error
}

/*

Entity is a folded Stream of <p,o> components

	entity := hexagon.Entity{}
  hexagon.
    Match(store, hexagon.IRI.Eq("s"), nil, nil).
	  FMap(entity.Append)

*/
type Entity map[curie.IRI]any

func (node Entity) Append(s, p curie.IRI, o any) error {
	if val, has := node[p]; !has {
		node[p] = o
	} else {
		switch v := val.(type) {
		case []any:
			node[p] = append(v, o)
		default:
			node[p] = []any{v, o}
		}
	}
	return nil
}

/*

Graph is a folded Stream of <s,p,o> components

	graph := hexagon.Graph{}
  hexagon.
    Match(store, nil, nil, hexagon.Eq("o")).
	  FMap(graph.Append)

*/
type Graph map[curie.IRI]Entity

func (graph Graph) Append(s, p curie.IRI, o any) error {
	if val, has := graph[s]; !has {
		node := Entity{}
		node.Append(s, p, o)
		graph[s] = node
	} else {
		val.Append(s, p, o)
		graph[s] = val
	}
	return nil
}

//
//
// Internal type
//
//

// components of <s,p,o,k> triple
type s = curie.IRI
type p = curie.IRI
type o = any
type k = guid.LID

// index types for 3rd faction
type __s = *skiplist.SkipList[s, k]
type __p = *skiplist.SkipList[p, k]
type __o = *skiplist.SkipList[o, k]

// index types for 2nd faction
type _po = *skiplist.SkipList[p, __o]
type _op = *skiplist.SkipList[o, __p]
type _so = *skiplist.SkipList[s, __o]
type _os = *skiplist.SkipList[o, __s]
type _sp = *skiplist.SkipList[s, __p]
type _ps = *skiplist.SkipList[p, __s]

// triple indexes
type spo = *skiplist.SkipList[s, _po]
type sop = *skiplist.SkipList[s, _op]
type pso = *skiplist.SkipList[p, _so]
type pos = *skiplist.SkipList[p, _os]
type osp = *skiplist.SkipList[o, _sp]
type ops = *skiplist.SkipList[o, _ps]

// allocators for indexes
func newS(rnd rand.Source) __s { return skiplist.New[s, k](ord.IRI, rnd) }
func newP(rnd rand.Source) __p { return skiplist.New[p, k](ord.IRI, rnd) }
func newO(rnd rand.Source) __o { return skiplist.New[o, k](ord.Any, rnd) }

func newPO(rnd rand.Source) _po { return skiplist.New[p, __o](ord.IRI, rnd) }
func newOP(rnd rand.Source) _op { return skiplist.New[o, __p](ord.Any, rnd) }
func newSO(rnd rand.Source) _so { return skiplist.New[s, __o](ord.IRI, rnd) }
func newOS(rnd rand.Source) _os { return skiplist.New[o, __s](ord.Any, rnd) }
func newSP(rnd rand.Source) _sp { return skiplist.New[s, __p](ord.IRI, rnd) }
func newPS(rnd rand.Source) _ps { return skiplist.New[p, __s](ord.IRI, rnd) }

func newSPO(rnd rand.Source) spo { return skiplist.New[s, _po](ord.IRI, rnd) }
func newSOP(rnd rand.Source) sop { return skiplist.New[s, _op](ord.IRI, rnd) }
func newPSO(rnd rand.Source) pso { return skiplist.New[p, _so](ord.IRI, rnd) }
func newPOS(rnd rand.Source) pos { return skiplist.New[p, _os](ord.IRI, rnd) }
func newOSP(rnd rand.Source) osp { return skiplist.New[o, _sp](ord.Any, rnd) }
func newOPS(rnd rand.Source) ops { return skiplist.New[o, _ps](ord.Any, rnd) }
