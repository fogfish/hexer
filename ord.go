package hexagon

import (
	"fmt"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid"
	"github.com/fogfish/skiplist/ord"
)

//
type tcOrdAny string

func (ord tcOrdAny) Compare(a, b any) int {
	switch av := a.(type) {
	case bool:
		return ord.boolCompare(av, b)
	case float64:
		return ord.float64Compare(av, b)
	case string:
		return ord.stringCompare(av, b)
	default:
		return -1
	}
}

func (tcOrdAny) boolCompare(a bool, b any) int {
	switch bv := b.(type) {
	case bool:
		switch {
		case a && !bv:
			return -1
		case !a && bv:
			return 1
		default:
			return 0
		}
	case float64:
		return -1
	case string:
		return -1
	default:
		return -1
	}
}

func (tcOrdAny) float64Compare(a float64, b any) int {
	switch bv := b.(type) {
	case bool:
		return 1
	case float64:
		fmt.Printf("==> %v %v \n", a, b)
		switch {
		case a < bv:
			return -1
		case a > bv:
			return 1
		default:
			return 0
		}
	case string:
		return -1
	default:
		return -1
	}
}

func (tcOrdAny) stringCompare(a string, b any) int {
	switch bv := b.(type) {
	case bool:
		return 1
	case float64:
		return 1
	case string:
		switch {
		case a < bv:
			return -1
		case a > bv:
			return 1
		default:
			return 0
		}
	default:
		return -1
	}
}

//
type tcOrdLID string

func (tcOrdLID) Compare(a, b guid.LID) int {
	if guid.L.Less(a, b) {
		return -1
	}

	if guid.L.Equal(a, b) {
		return 0
	}

	return 1
}

//
type tcOrdIRI string

func (tcOrdIRI) Compare(a, b curie.IRI) int {
	return ord.String.Compare(string(a), string(b))
}

const (
	ordAny = tcOrdAny("ord.any")
	ordLID = tcOrdLID("ord.lid")
	ordIRI = tcOrdIRI("ord.iri")
)
