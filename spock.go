package hexer

import (
	"github.com/fogfish/curie"
	"github.com/fogfish/guid/v2"
)

const (
	XSD_ANYURI = curie.IRI("xsd:anyURI")
	XSD_STRING = curie.IRI("xsd:string")
)

type XSDValue interface{ isXSDValue() }

type XSDAnyURI struct{ Value curie.IRI }

func (XSDAnyURI) isXSDValue() {}

type XSDString struct{ Value string }

func (XSDString) isXSDValue() {}

type SPOCK struct {
	S, P curie.IRI
	O    XSDValue
	C    float64
	K    guid.K
}

// func (SPOCK[T]) HKT1(Type) {}
// func (SPOCK[T]) HKT2(T)    {}
