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

func ToXSDValue[T DataType](value T) XSDValue {
	switch v := any(value).(type) {
	case string:
		return XSDString{Value: v}
	default:
		panic("xxxx")
	}
}

type SPOCK struct {
	S curie.IRI
	P curie.IRI
	O XSDValue
	C float64
	K guid.K
}

type Pattern struct {
	S *Predicate[curie.IRI]
	P *Predicate[curie.IRI]
	O *Predicate[XSDValue]
}

func Query(s *Predicate[curie.IRI], p *Predicate[curie.IRI], o *Predicate[XSDValue]) Pattern {
	return Pattern{S: s, P: p, O: o}
}

func (q Pattern) Strategy() Strategy {
	switch {
	// x, o, _ ⇒ spo
	case exact(q.S) && order(q.P) && !exact(q.O):
		return 0510
	// x, _, o ⇒ sop
	case exact(q.S) && !exact(q.P) && order(q.O):
		return 0501
	// _, x, o ⇒ pos
	case !exact(q.S) && exact(q.P) && order(q.O):
		return 0051
	// o, x, _ ⇒ pso
	case order(q.S) && exact(q.P) && !exact(q.O):
		return 0150
	// o, _, x ⇒ osp
	case order(q.S) && !exact(q.P) && exact(q.O):
		return 0105
	// _, o, x ⇒ ops
	case !exact(q.S) && order(q.P) && exact(q.O):
		return 0015

	// x, x, _ ⇒ spo
	case exact(q.S) && exact(q.P) && !exact(q.O):
		return 0550
	// _, x, x ⇒ pos
	case !exact(q.S) && exact(q.P) && exact(q.O):
		return 0055
	// x, _, x ⇒ sop
	case exact(q.S) && !exact(q.P) && exact(q.O):
		return 0505

	// x, _, _ ⇒ spo
	case exact(q.S) && !exact(q.P) && !exact(q.O):
		return 0500
	// _, x, _ ⇒ pso
	case !exact(q.S) && exact(q.P) && !exact(q.O):
		return 0050
	// _, _, x ⇒ osp
	case !exact(q.S) && !exact(q.P) && exact(q.O):
		return 0005

	// _, _, _ ⇒ spo
	case !exact(q.S) && !exact(q.P) && !exact(q.O):
		return 0000
	}

	return 0777
}

type StreamX interface {
	Head() SPOCK
	Next() bool
	// FMap(f func(SPOCK) error) error
}
