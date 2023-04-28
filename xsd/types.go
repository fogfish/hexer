package xsd

import (
	"strconv"

	"github.com/fogfish/curie"
)

type Value interface{ XSDType() curie.IRI }

// The data type represents Internationalized Resource Identifier.
// Used to uniquely identify concept, objects, etc.
type AnyURI curie.IRI

const XSD_ANYURI = curie.IRI("xsd:anyURI")

func (v AnyURI) XSDType() curie.IRI { return XSD_ANYURI }
func (v AnyURI) String() string     { return curie.IRI(v).Safe() }

// The string data-type represents character strings in knowledge statements.
// The language strings are annotated with corresponding language tag.
type String string

const XSD_STRING = curie.IRI("xsd:string")

func (v String) XSDType() curie.IRI { return XSD_STRING }
func (v String) String() string     { return strconv.Quote(string(v)) }

// The Integer data-type in knowledge statement.
// The library uses various int precision data-types to represent decimal values.
// type XSDInteger = int
// type Byte = int8
// type Short = int16
// type Int = int32
// type Long = int64
// type NonNegativeInteger = uint
// type UnsignedByte = uint8
// type UnsignedShort = uint16
// type UnsignedInt = uint32
// type UnsignedLong = uint64

// const XSD_INTEGER = curie.IRI("xsd:integer")

// type Integer struct{ Value XSDInteger }

// func (Integer) XSDType() curie.IRI { return XSD_INTEGER }
