package hexer

import (
	"strings"

	"github.com/fogfish/curie"
	"github.com/fogfish/hexer/xsd"
)

// Stream of knowledge statements ⟨s,p,o,c,k⟩
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

func NewFilterO(hint Hint, q *Predicate[xsd.Value], stream Stream) Stream {
	switch hint {
	case HINT_MATCH:
		return NewFilter(
			func(spock SPOCK) bool { return xsd.Compare(spock.O, q.Value) == 0 },
			stream,
		)
	case HINT_FILTER_PREFIX:
		return NewFilter(
			func(spock SPOCK) bool { return xsd.HasPrefix(spock.O, q.Value) },
			stream,
		)
	case HINT_FILTER:
		switch q.Clause {
		case LT:
			return NewFilter(
				func(spock SPOCK) bool { return xsd.Compare(spock.O, q.Value) == -1 },
				stream,
			)
		case GT:
			return NewFilter(
				func(spock SPOCK) bool { return xsd.Compare(spock.O, q.Value) == 1 },
				stream,
			)
		case IN:
			return NewFilter(
				func(spock SPOCK) bool {
					return xsd.Compare(spock.O, q.Value) >= 0 && xsd.Compare(spock.O, q.Other) <= 0
				},
				stream,
			)
		}
	}

	return stream
}

func NewFilterP(hint Hint, q *Predicate[curie.IRI], stream Stream) Stream {
	switch hint {
	case HINT_MATCH:
		return NewFilter(
			func(spock SPOCK) bool { return spock.P == q.Value },
			stream,
		)
	case HINT_FILTER_PREFIX:
		return NewFilter(
			func(spock SPOCK) bool { return strings.HasPrefix(string(spock.P), string(q.Value)) },
			stream,
		)
	}

	return stream
}

func NewFilterS(hint Hint, q *Predicate[curie.IRI], stream Stream) Stream {
	switch hint {
	case HINT_MATCH:
		return NewFilter(
			func(spock SPOCK) bool { return spock.S == q.Value },
			stream,
		)
	case HINT_FILTER_PREFIX:
		return NewFilter(
			func(spock SPOCK) bool { return strings.HasPrefix(string(spock.S), string(q.Value)) },
			stream,
		)
	}

	return stream
}
