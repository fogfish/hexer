package hexer

import (
	"fmt"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid/v2"
	"github.com/fogfish/hexer/xsd"
)

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

func From[T xsd.DataType](s, p curie.IRI, o T) SPOCK {
	return SPOCK{S: s, P: p, O: xsd.From(o)}
}

type Bag []SPOCK

func (bag *Bag) Join(spock SPOCK) error {
	*bag = append(*bag, spock)
	return nil
}

// Stream of <s,p,o,c,k> triples fetched from the store
type Stream interface {
	Head() SPOCK
	Next() bool
	FMap(func(SPOCK) error) error
}

type filter struct {
	pred   func(SPOCK) bool
	stream Stream
}

func (filter *filter) Head() SPOCK {
	return filter.stream.Head()
}

func (filter *filter) Next() bool {
	for {
		if !filter.stream.Next() {
			return false
		}

		if filter.pred(filter.stream.Head()) {
			return true
		}
	}
}

func (filter *filter) FMap(f func(SPOCK) error) error {
	for filter.Next() {
		if err := f(filter.Head()); err != nil {
			return err
		}
	}
	return nil
}

func NewFilter(pred func(SPOCK) bool, stream Stream) Stream {
	return &filter{pred: pred, stream: stream}
}
