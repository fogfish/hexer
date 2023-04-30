package xsd

import "github.com/fogfish/curie"

// DataType is a type constrain used by the library.
// See https://www.w3.org/TR/xmlschema-2/#datatype
//
// Knowledge statements contain scalar objects -- literals. Literals are either
// language-tagged string `rdf:langString` or type-safe values containing a
// reference to data-type (e.g. `xsd:string`).
//
// This interface defines data-types supported by the library. It maps well-known
// semantic types to Golang native types and relation to existed schema(s) and
// ontologies.
type DataType interface {
	~string |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~bool |
		~[]byte
}

// The floating point data-type in knowledge statement.
// The library uses various uint precisions.
type Float = float32
type Double = float64

// The boolean data-type in knowledge statement
type Boolean = bool

type HexBinary = []byte
type Base64Binary = []byte

// From builds Object from Golang type
func From[T DataType](value T) Value {
	switch v := any(value).(type) {
	case curie.IRI:
		return AnyURI(v)
	case AnyURI:
		return v
	case string:
		return String(v)
	case String:
		return v
	// case int:
	// 	return Integer{Value: v}
	default:
		panic("xxxx")
	}
}
