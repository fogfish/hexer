package hexer

import (
	"github.com/fogfish/curie"
	"github.com/fogfish/hexer/xsd"
)

// types of predicate clauses
type Clause int

const (
	ALL Clause = iota
	EQ         // Equal
	PQ         // Prefix Equal
	LT         // Less Than
	GT         // Greater Than
	IN         // InRange, Between
)

// Predicate on <s,p,o>
type Predicate[T any] struct {
	Clause Clause
	Value  T
	Other  T
}

// Makes `equal` to IRI predicate
func IRI(value curie.IRI) *Predicate[curie.IRI] {
	return &Predicate[curie.IRI]{Clause: EQ, Value: value}
}

// Makes `equal` to IRI predicate
func Prefix(value curie.IRI) *Predicate[curie.IRI] {
	return &Predicate[curie.IRI]{Clause: PQ, Value: value}
}

// Makes `equal to` value predicate
func Eq[T xsd.DataType](value T) *Predicate[Object] {
	return &Predicate[Object]{Clause: EQ, Value: From(value)}
}
