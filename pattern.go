//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/hexagon
//

package hexer

import (
	"github.com/fogfish/curie"
	"github.com/fogfish/skiplist"
)

// types of predicate clauses
type Clause int

const (
	ALL Clause = iota
	EQ
	LT
	GT
	IN
)

// Predicate on <s,p,o>
type Predicate[T any] struct {
	Clause Clause
	Value  T
	other  T
}

type iri string

const IRI = iri("hexagon.iri")

// Makes `equal` to IRI predicate
func (iri) Eq(value curie.IRI) *Predicate[curie.IRI] {
	return &Predicate[curie.IRI]{Clause: EQ, Value: value}
}

// Makes `less than` IRI predicate
func (iri) Lt(value curie.IRI) *Predicate[curie.IRI] {
	return &Predicate[curie.IRI]{Clause: LT, Value: value}
}

// Makes `greater than` IRI predicate
func (iri) Gt(value curie.IRI) *Predicate[curie.IRI] {
	return &Predicate[curie.IRI]{Clause: GT, Value: value}
}

// Makes `equal to` value predicate
func Eq[T DataType](value T) *Predicate[XSDValue] {
	return &Predicate[XSDValue]{Clause: EQ, Value: ToXSDValue(value)}
}

// Makes `less than` value predicate
func Lt[T DataType](value T) *Predicate[XSDValue] {
	return &Predicate[XSDValue]{Clause: LT, Value: ToXSDValue(value)}
}

// Makes `greater than` value predicate
func Gt[T DataType](value T) *Predicate[XSDValue] {
	return &Predicate[XSDValue]{Clause: GT, Value: ToXSDValue(value)}
}

// Makes `in range` predicate
func In[T DataType](from, to T) *Predicate[XSDValue] {
	return &Predicate[XSDValue]{Clause: IN, Value: ToXSDValue(from), other: ToXSDValue(to)}
}

// checks if predicate is exact match
func exact[K any](pred *Predicate[K]) bool {
	return pred != nil && pred.Clause == EQ
}

// checks if predicate is order/filter
func order[K any](pred *Predicate[K]) bool {
	return pred != nil && pred.Clause != EQ
}

// toIterator builds Iterator from skip list using predicate
func toIterator[A, B any](pred *Predicate[A], list *skiplist.SkipList[A, B]) *skiplist.Iterator[A, B] {
	switch {
	case pred == nil:
		return skiplist.Values(list)
	case pred.Clause == LT:
		before, _ := skiplist.Split(list, pred.Value)
		return before
	case pred.Clause == GT:
		_, after := skiplist.Split(list, pred.Value)
		return after
	case pred.Clause == IN:
		return skiplist.Range(list, pred.Value, pred.other)
	}

	return nil
}

/*
Strategy is a code that defines combination of indexes and resolution
strategy to be used for query. The code consists of octal digits for
each index in the order of <s,p,o>:
  - 0: best effort lookup, storage tries to scope lookup with filters
  - 1: lookup uses filters defined by predicate so that multiple values are inspected
  - 5: lookup uses exact by predicate so that single value is inspected
*/
type Strategy int

func (s Strategy) String() string {
	switch s {
	case 0510:
		return "(s)ᴾ ⇒ o"
	// x, _, o ⇒ sop
	case 0501:
		return "(s)º ⇒ p"
	// _, x, o ⇒ pos
	case 0051:
		return "(p)º ⇒ s"
	// o, x, _ ⇒ pso
	case 0150:
		return "(p)ˢ ⇒ o"
	// o, _, x ⇒ osp
	case 0105:
		return "(o)ˢ ⇒ p"
	// _, o, x ⇒ ops
	case 0015:
		return "(o)ᴾ ⇒ s"

	// x, x, _ ⇒ spo
	case 0550:
		return "(sp) ⇒ o"
	// _, x, x ⇒ pos
	case 0055:
		return "(po) ⇒ s"
	// x, _, x ⇒ sop
	case 0505:
		return "(so) ⇒ p"

	// x, _, _ ⇒ spo
	case 0500:
		return "(s) ⇒ po"
	// _, x, _ ⇒ pso
	case 0050:
		return "(p) ⇒ so"
	// _, _, x ⇒ osp
	case 0005:
		return "(o) ⇒ sp"

	// _, _, _ ⇒ spo
	case 0000:
		return "∅ ⇒ spo"
	}

	return ""
}

type pattern struct {
	store *Store
	s     *Predicate[s]
	p     *Predicate[p]
	o     *Predicate[o]
}

// evaluates pattern
func (q pattern) eval() (Strategy, Stream) {
	strategy := q.strategy()
	switch strategy {
	// x, o, _ ⇒ spo
	case 0510:
		return strategy, q.sPO()
	// x, _, o ⇒ sop
	case 0501:
		return strategy, q.sOP()
	// _, x, o ⇒ pos
	case 0051:
		return strategy, q.pOS()
	// o, x, _ ⇒ pso
	case 0150:
		return strategy, q.pSO()
	// o, _, x ⇒ osp
	case 0105:
		return strategy, q.oSP()
	// _, o, x ⇒ ops
	case 0015:
		return strategy, q.oPS()

	// x, x, _ ⇒ spo
	case 0550:
		return strategy, q.spO()
	// _, x, x ⇒ pos
	case 0055:
		return strategy, q.poS()
	// x, _, x ⇒ sop
	case 0505:
		return strategy, q.soP()

	// x, _, _ ⇒ spo
	case 0500:
		return strategy, q.sPO()
	// _, x, _ ⇒ pso
	case 0050:
		return strategy, q.pSO()
	// _, _, x ⇒ osp
	case 0005:
		return strategy, q.oSP()

	// _, _, _ ⇒ spo
	case 0000:
		return strategy, q.spo()
	}

	return strategy, nil
}

/*
builds execution strategy for the pattern

x, o, _ ⇒ spo
x, _, o ⇒ sop
_, x, o ⇒ pos
o, x, _ ⇒ pso
o, _, x ⇒ osp
_, o, x ⇒ ops

x, x, _ ⇒ spo
_, x, x ⇒ pos
x, _, x ⇒ sop

x, _, _ ⇒ spo
_, x, _ ⇒ pso
_, _, x ⇒ osp

_, _, _ ⇒ spo
*/
func (q pattern) strategy() Strategy {
	switch {
	// x, o, _ ⇒ spo
	case exact(q.s) && order(q.p) && !exact(q.o):
		return 0510
	// x, _, o ⇒ sop
	case exact(q.s) && !exact(q.p) && order(q.o):
		return 0501
	// _, x, o ⇒ pos
	case !exact(q.s) && exact(q.p) && order(q.o):
		return 0051
	// o, x, _ ⇒ pso
	case order(q.s) && exact(q.p) && !exact(q.o):
		return 0150
	// o, _, x ⇒ osp
	case order(q.s) && !exact(q.p) && exact(q.o):
		return 0105
	// _, o, x ⇒ ops
	case !exact(q.s) && order(q.p) && exact(q.o):
		return 0015

	// x, x, _ ⇒ spo
	case exact(q.s) && exact(q.p) && !exact(q.o):
		return 0550
	// _, x, x ⇒ pos
	case !exact(q.s) && exact(q.p) && exact(q.o):
		return 0055
	// x, _, x ⇒ sop
	case exact(q.s) && !exact(q.p) && exact(q.o):
		return 0505

	// x, _, _ ⇒ spo
	case exact(q.s) && !exact(q.p) && !exact(q.o):
		return 0500
	// _, x, _ ⇒ pso
	case !exact(q.s) && exact(q.p) && !exact(q.o):
		return 0050
	// _, _, x ⇒ osp
	case !exact(q.s) && !exact(q.p) && exact(q.o):
		return 0005

	// _, _, _ ⇒ spo
	case !exact(q.s) && !exact(q.p) && !exact(q.o):
		return 0000
	}

	return 0777
}

// ∅ ⇒ spo
func (q pattern) spo() Stream {
	iter := &iterator3{spo: skiplist.Values(q.store.spo)}

	return iter
}

// (s)ᴾ ⇒ o
func (q pattern) sPO() Stream {
	iter := &iterator2[p, o]{s: q.s.Value, ap: q.o}

	if _po, has := skiplist.Lookup(q.store.spo, q.s.Value); has {
		iter.a = &iter.o
		iter.b = &iter.p
		iter._ba = toIterator(q.p, _po)

		return iter
	}

	return iter
}

// (s)º ⇒ p
func (q pattern) sOP() Stream {
	iter := &iterator2[o, p]{s: q.s.Value, ap: q.p}

	if _op, has := skiplist.Lookup(q.store.sop, q.s.Value); has {
		iter.a = &iter.p
		iter.b = &iter.o
		iter._ba = toIterator(q.o, _op)

		return iter
	}

	return iter
}

// (p)ˢ ⇒ o
func (q pattern) pSO() Stream {
	iter := &iterator2[s, o]{p: q.p.Value, ap: q.o}

	if _so, has := skiplist.Lookup(q.store.pso, q.p.Value); has {
		iter.a = &iter.o
		iter.b = &iter.s
		iter._ba = toIterator(q.s, _so)

		return iter
	}

	return iter
}

// (p)º ⇒ s
func (q pattern) pOS() Stream {
	iter := &iterator2[o, s]{p: q.p.Value, ap: q.s}

	if _os, has := skiplist.Lookup(q.store.pos, q.p.Value); has {
		iter.a = &iter.s
		iter.b = &iter.o
		iter._ba = toIterator(q.o, _os)

		return iter
	}

	return iter
}

// (o)ˢ ⇒ p
func (q pattern) oSP() Stream {
	iter := &iterator2[s, p]{o: q.o.Value, ap: q.p}

	if _sp, has := skiplist.Lookup(q.store.osp, q.o.Value); has {
		// fmt.Printf("==> len o[sp] %v\n", skiplist.Length(_sp))
		iter.a = &iter.p
		iter.b = &iter.s
		iter._ba = toIterator(q.s, _sp)

		return iter
	}

	return iter
}

// (o)ᴾ ⇒ s
func (q pattern) oPS() Stream {
	iter := &iterator2[p, s]{o: q.o.Value, ap: q.s}

	if _ps, has := skiplist.Lookup(q.store.ops, q.o.Value); has {
		iter.a = &iter.s
		iter.b = &iter.p
		iter._ba = toIterator(q.p, _ps)

		return iter
	}

	return iter
}

// (sp) ⇒ o
func (q pattern) spO() Stream {
	iter := &iterator[o]{s: q.s.Value, p: q.p.Value}

	if _po, has := skiplist.Lookup(q.store.spo, q.s.Value); has {
		if __o, has := skiplist.Lookup(_po, q.p.Value); has {
			iter.a = &iter.o
			iter.__a = toIterator(q.o, __o)

			return iter
		}
	}

	return iter
}

// (po) ⇒ s
func (q pattern) poS() Stream {
	iter := &iterator[s]{p: q.p.Value, o: q.o.Value}

	if _os, has := skiplist.Lookup(q.store.pos, q.p.Value); has {
		if __s, has := skiplist.Lookup(_os, q.o.Value); has {
			iter.a = &iter.s
			iter.__a = toIterator(q.s, __s)

			return iter
		}
	}

	return iter
}

// (so) ⇒ p
func (q pattern) soP() Stream {
	iter := &iterator[p]{s: q.s.Value, o: q.o.Value}

	if _op, has := skiplist.Lookup(q.store.sop, q.s.Value); has {
		if __p, has := skiplist.Lookup(_op, q.o.Value); has {
			iter.a = &iter.p
			iter.__a = toIterator(q.p, __p)

			return iter
		}
	}

	return iter
}
