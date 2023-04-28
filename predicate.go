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

// TODO: consolidate interface for query definition
// When to use
//  - IRI
//  - Like
//  - Ref
//  - Eq
//  - Prefix
//

// TODO: Prefix only for curie.IRI & xsd.String

// Makes `equal` to IRI predicate
func IRI(value curie.IRI) *Predicate[curie.IRI] {
	return &Predicate[curie.IRI]{Clause: EQ, Value: value}
}

// Makes `equal` to IRI predicate
func Like(value curie.IRI) *Predicate[curie.IRI] {
	return &Predicate[curie.IRI]{Clause: PQ, Value: value}
}

// Makes `equal to` value predicate
func Ref(value curie.IRI) *Predicate[xsd.Value] {
	return &Predicate[xsd.Value]{Clause: EQ, Value: xsd.AnyURI(value)}
}

// Makes `equal to` value predicate
func Eq[T xsd.DataType](value T) *Predicate[xsd.Value] {
	return &Predicate[xsd.Value]{Clause: EQ, Value: xsd.From(value)}
}

// Makes `prefix` value predicate
func Prefix[T xsd.DataType](value T) *Predicate[xsd.Value] {
	return &Predicate[xsd.Value]{Clause: PQ, Value: xsd.From(value)}
}

// Makes `less than` value predicate
func Lt[T xsd.DataType](value T) *Predicate[xsd.Value] {
	return &Predicate[xsd.Value]{Clause: LT, Value: xsd.From(value)}
}

// Makes `greater than` value predicate
func Gt[T xsd.DataType](value T) *Predicate[xsd.Value] {
	return &Predicate[xsd.Value]{Clause: GT, Value: xsd.From(value)}
}

// Makes `in range` predicate
func In[T xsd.DataType](from, to T) *Predicate[xsd.Value] {
	return &Predicate[xsd.Value]{Clause: IN, Value: xsd.From(from), Other: xsd.From(to)}
}
