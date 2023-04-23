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

// The data type represents Internationalized Resource Identifier.
// Used to uniquely identify concept, objects, etc.
type AnyURI = curie.IRI

// The string data-type represents character strings in knowledge statements.
// The language strings are annotated with corresponding language tag.
type String = string

// The Integer data-type in knowledge statement.
// The library uses various int precision data-types to represent decimal values.
type Integer = int
type Byte = int8
type Short = int16
type Int = int32
type Long = int64
type NonNegativeInteger = uint
type UnsignedByte = uint8
type UnsignedShort = uint16
type UnsignedInt = uint32
type UnsignedLong = uint64

// The floating point data-type in knowledge statement.
// The library uses various uint precisions.
type Float = float32
type Double = float64

// The boolean data-type in knowledge statement
type Boolean = bool

type HexBinary = []byte
type Base64Binary = []byte

// const (
// 	ANYURI  = curie.IRI("xsd:anyURI")
// 	STRING  = curie.IRI("xsd:string")
// 	INTEGER = curie.IRI("xsd:integer")
// 	BYTE    = curie.IRI("xsd:byte")
// 	SHORT   = curie.IRI("xsd:short")
// 	INT     = curie.IRI("xsd:int")
// 	LONG    = curie.IRI("xsd:long")
// )
