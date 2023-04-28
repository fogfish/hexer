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

func Link(s, p, o curie.IRI) SPOCK {
	return SPOCK{S: s, P: p, O: xsd.AnyURI(o)}
}

func From[T xsd.DataType](s, p curie.IRI, o T) SPOCK {
	return SPOCK{S: s, P: p, O: xsd.From(o)}
}

type Bag []SPOCK

func (bag *Bag) Ref(s, p, o curie.IRI) {
	*bag = append(*bag, SPOCK{S: s, P: p, O: xsd.AnyURI(o)})
}

func (bag *Bag) Add(s, p curie.IRI, o xsd.Value) {
	*bag = append(*bag, SPOCK{S: s, P: p, O: o})
}

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

/*

	h := filter.seq.Head()
	switch {
	case filter.q.Clause == hexer.LT && h.S < filter.q.Value:
		return true
	case filter.q.Clause == hexer.GT && h.S > filter.q.Value:
		return true
	case filter.q.Clause == hexer.IN && h.S > filter.q.Value && h.S < filter.q.Other:
		return true
	}


*/

// type FilterP struct {
// 	seq hexer.Stream
// 	q   *hexer.Predicate[curie.IRI]
// }

// func (filter *FilterP) Head() hexer.SPOCK {
// 	return filter.seq.Head()
// }

// func (filter *FilterP) Next() bool {
// 	for {
// 		if !filter.seq.Next() {
// 			return false
// 		}
// 		h := filter.seq.Head()
// 		switch {
// 		case filter.q.Clause == hexer.LT && h.P < filter.q.Value:
// 			return true
// 		case filter.q.Clause == hexer.GT && h.P > filter.q.Value:
// 			return true
// 		case filter.q.Clause == hexer.IN && h.P > filter.q.Value && h.S < filter.q.Other:
// 			return true
// 		}
// 	}
// }

// type FilterO struct {
// 	seq hexer.StreamX
// 	q   *hexer.Predicate[hexer.XSDValue]
// }
//
// func (filter FilterO) Head() hexer.SPOCK {
// 	return filter.seq.Head()
// }
//
// func (filter FilterO) Next() bool {
// 	for {
// 		if !filter.seq.Next() {
// 			return false
// 		}
// 		h := filter.seq.Head()
// 		switch {
// 		case filter.q.Clause == hexer.LT && h.O < filter.q.Value:
// 			return true
// 		case filter.q.Clause == hexer.GT && h.O > filter.q.Value:
// 			return true
// 		case filter.q.Clause == hexer.IN && h.O > filter.q.Value && h.O < filter.q.Other:
// 			return true
// 		}
// 	}
// }
