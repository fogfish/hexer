package hexer

import (
	"fmt"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid/v2"
	"github.com/fogfish/hexer/xsd"
)

// Knowledge statement
//
//	s: subject
//	p: predicate
//	o: object
//	c: credibility
//	k: k-order
type SPOCK struct {
	S curie.IRI
	P curie.IRI
	O xsd.Value
	C float64
	K guid.K
}

func (spock SPOCK) String() string {
	return fmt.Sprintf("⟨%s %s %s⟩", spock.S.Safe(), spock.P.Safe(), spock.O)
}

// Create new knowledge statement From
func From[T xsd.DataType](s, p curie.IRI, o T) SPOCK {
	return SPOCK{S: s, P: p, O: xsd.From(o)}
}

// Collection of knowledge statements
type Bag []SPOCK

func (bag *Bag) Join(spock SPOCK) error {
	*bag = append(*bag, spock)
	return nil
}
