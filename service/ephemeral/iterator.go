package ephemeral

import (
	"strings"

	"github.com/fogfish/curie"
	"github.com/fogfish/hexer"
	"github.com/fogfish/hexer/xsd"
	"github.com/fogfish/skiplist"
)

// TODO: use xsd as primary type for indexes

type Seq[K, V any] interface {
	Head() (K, V)
	Next() bool
}

func overIRI[A, B any](
	pred *hexer.Predicate[curie.IRI],
) func(list *skiplist.SkipList[curie.IRI, B]) Seq[A, B] {
	return func(list *skiplist.SkipList[curie.IRI, B]) Seq[A, B] {
		var seq Seq[curie.IRI, B]

		switch {
		case pred == nil:
			seq = skiplist.Values(list)
		case pred.Clause == hexer.EQ:
			seq = skiplist.Slice(list, pred.Value, 1)
		case pred.Clause == hexer.PQ:
			_, after := skiplist.Split(list, pred.Value)
			seq = after.TakeWhile(
				func(x curie.IRI) bool {
					return strings.HasPrefix(string(x), string(pred.Value))
				},
			)
			// TODO: block after prefix exceeds
		case pred.Clause == hexer.LT:
			seq, _ = skiplist.Split(list, pred.Value)
		case pred.Clause == hexer.GT:
			_, seq = skiplist.Split(list, pred.Value)
		case pred.Clause == hexer.IN:
			seq = skiplist.Range(list, pred.Value, pred.Other)
		}

		return seq.(Seq[A, B])
	}
}

func overXSD[A, B any](
	pred *hexer.Predicate[xsd.Value],
) func(list *skiplist.SkipList[xsd.Value, B]) Seq[A, B] {
	return func(list *skiplist.SkipList[xsd.Value, B]) Seq[A, B] {
		var seq Seq[xsd.Value, B]

		switch {
		case pred == nil:
			seq = skiplist.Values(list)
		case pred.Clause == hexer.EQ:
			seq = skiplist.Slice(list, pred.Value, 1)
		case pred.Clause == hexer.PQ:
			_, after := skiplist.Split(list, pred.Value)
			seq = after.TakeWhile(
				func(x xsd.Value) bool {
					return xsd.HasPrefix(x, pred.Value)
				},
			)
			// TODO: block after prefix exceeds
		case pred.Clause == hexer.LT:
			seq, _ = skiplist.Split(list, pred.Value)
		case pred.Clause == hexer.GT:
			_, seq = skiplist.Split(list, pred.Value)
		case pred.Clause == hexer.IN:
			seq = skiplist.Range(list, pred.Value, pred.Other)
		}

		return seq.(Seq[A, B])
	}
}

func toIterator[A, B any](
	pred *hexer.Predicate[A],
	list *skiplist.SkipList[A, B],
) *skiplist.Iterator[A, B] {
	switch {
	case pred == nil:
		return skiplist.Values(list)
	case pred.Clause == hexer.EQ:
		return skiplist.Slice(list, pred.Value, 1)
	case pred.Clause == hexer.PQ:
		switch v := any(pred.Value).(type) {
		case curie.IRI:
			_, after := skiplist.Split(list, pred.Value)
			return after.TakeWhile(func(a A) bool {
				if iri, ok := any(a).(curie.IRI); ok {
					return strings.HasPrefix(string(iri), string(v))
				}
				return false
			})
		default:
			panic("xxx")
		}

		// TODO: block after prefix exceeds
	case pred.Clause == hexer.LT:
		before, _ := skiplist.Split(list, pred.Value)
		return before
	case pred.Clause == hexer.GT:
		_, after := skiplist.Split(list, pred.Value)
		return after
	case pred.Clause == hexer.IN:
		return skiplist.Range(list, pred.Value, pred.Other)
	}

	return nil
}

type Iterator[A, B, C any] struct {
	a     A
	b     B
	pb    *hexer.Predicate[B]
	pc    *hexer.Predicate[C]
	abc   Seq[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]
	_bc   Seq[B, *skiplist.SkipList[C, k]]
	__c   Seq[C, k]
	f     func(A, B, C) (hexer.SPOCK, bool)
	spock hexer.SPOCK
	hlp   Helper[A, B, C]
}

type Helper[A, B, C any] interface {
	L1(*skiplist.SkipList[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]) Seq[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]
	L2(*skiplist.SkipList[B, *skiplist.SkipList[C, k]]) Seq[B, *skiplist.SkipList[C, k]]
	L3(*skiplist.SkipList[C, k]) Seq[C, k]
}

func NewIterator[A, B, C any](
	hlp Helper[A, B, C],
	seq *skiplist.SkipList[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]],
	// pa *hexer.Predicate[A],
	// pb *hexer.Predicate[B],
	// pc *hexer.Predicate[C],
	f func(A, B, C) (hexer.SPOCK, bool),
) *Iterator[A, B, C] {
	return &Iterator[A, B, C]{
		hlp: hlp,
		abc: hlp.L1(seq), //toIterator(pa, seq),
		// pb:  pb,
		// pc:  pc,
		f: f,
	}
}

func (iter *Iterator[A, B, C]) Head() hexer.SPOCK {
	return iter.spock
}

func (iter *Iterator[A, B, C]) Next() bool {
	if iter._bc == nil {
		if iter.abc == nil || !iter.abc.Next() {
			return false
		}
		a, _bc := iter.abc.Head()
		iter.a = a
		iter._bc = iter.hlp.L2(_bc) // toIterator(iter.pb, _bc)
	}

	if iter.__c == nil {
		if iter._bc == nil || !iter._bc.Next() {
			iter._bc = nil
			return iter.Next()
		}

		b, __c := iter._bc.Head()
		iter.b = b
		iter.__c = iter.hlp.L3(__c) // toIterator(iter.pc, __c)
	}

	if iter.__c == nil || !iter.__c.Next() {
		iter.__c = nil
		return iter.Next()
	}

	c, _ := iter.__c.Head()
	spock, has := iter.f(iter.a, iter.b, c)
	if !has {
		return false
	}
	iter.spock = spock

	return true
}

func (iter *Iterator[A, B, C]) FMap(f func(hexer.SPOCK) error) error {
	for iter.Next() {
		if err := f(iter.Head()); err != nil {
			return err
		}
	}
	return nil
}

// type BiIterator[B, A any] struct {
// 	s   s
// 	p   p
// 	o   o
// 	a   *A
// 	b   *B
// 	ap  *hexer.Predicate[A]
// 	_ba *skiplist.Iterator[B, *skiplist.SkipList[A, k]]
// 	__a *skiplist.Iterator[A, k]
// }

// func (iter *BiIterator[A, B]) Head() hexer.SPOCK {
// 	*iter.a, _ = iter.__a.Head()
// 	return hexer.SPOCK{S: iter.s, P: iter.p, O: iter.o}
// }

// func (iter *BiIterator[A, B]) Next() bool {
// 	if iter.__a == nil {
// 		if iter._ba == nil || !iter._ba.Next() {
// 			return false
// 		}
// 		b, __a := iter._ba.Head()
// 		*iter.b = b
// 		iter.__a = toIterator(iter.ap, __a)
// 	}

// 	if iter.__a == nil || !iter.__a.Next() {
// 		iter.__a = nil
// 		return iter.Next()
// 	}

// 	return true
// }

// func (iter *BiIterator[A, B]) FMap(f func(hexer.SPOCK) error) error {
// 	for iter.Next() {
// 		if err := f(iter.Head()); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// type Iterator[A any] struct {
// 	s   s
// 	p   p
// 	o   o
// 	a   *A
// 	__a *skiplist.Iterator[A, k]
// }

// func (iter *Iterator[A]) Head() hexer.SPOCK {
// 	*iter.a, _ = iter.__a.Head()
// 	return hexer.SPOCK{S: iter.s, P: iter.p, O: iter.o}
// }

// func (iter *Iterator[A]) Next() bool {
// 	return iter.__a != nil && iter.__a.Next()
// }

// func (iter *Iterator[T]) FMap(f func(hexer.SPOCK) error) error {
// 	for iter.Next() {
// 		if err := f(iter.Head()); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// type TriIterator struct {
// 	s   s
// 	p   p
// 	spo *skiplist.Iterator[s, _po]
// 	_po *skiplist.Iterator[p, __o]
// 	__o *skiplist.Iterator[o, k]
// }

// func (iter *TriIterator) Head() hexer.SPOCK {
// 	o, k := iter.__o.Head()

// 	return hexer.SPOCK{S: iter.s, P: iter.p, O: o, K: k}
// }

// func (iter *TriIterator) Next() bool {
// 	if iter._po == nil {
// 		if iter.spo == nil || !iter.spo.Next() {
// 			return false
// 		}
// 		s, _po := iter.spo.Head()
// 		iter.s = s
// 		iter._po = skiplist.Values(_po)
// 	}

// 	if iter.__o == nil {
// 		if iter._po == nil || !iter._po.Next() {
// 			iter._po = nil
// 			return iter.Next()
// 		}

// 		p, __o := iter._po.Head()
// 		iter.p = p
// 		iter.__o = skiplist.Values(__o)
// 	}

// 	if iter.__o == nil || !iter.__o.Next() {
// 		iter.__o = nil
// 		return iter.Next()
// 	}

// 	return true
// }

// func (iter *TriIterator) FMap(f func(hexer.SPOCK) error) error {
// 	for iter.Next() {
// 		if err := f(iter.Head()); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// TODO:
//  - combinator for itterators (generic one)
//  - join function for iterators
//  1. generic iterator with predicate for __x
//  2. generic iterator with predicate for _x_ & x__
//

// type iterator2[B, A any] struct {
// 	s   s
// 	p   p
// 	o   o
// 	a   *A
// 	b   *B
// 	ap  *Predicate[A]
// 	_ba *skiplist.Iterator[B, *skiplist.SkipList[A, k]]
// 	__a *skiplist.Iterator[A, k]
// }

// func (iter *iterator2[A, B]) Head() (s, p, o) {
// 	*iter.a, _ = iter.__a.Head()
// 	return iter.s, iter.p, iter.o
// }

// func (iter *iterator2[A, B]) Next() bool {
// 	if iter.__a == nil {
// 		if iter._ba == nil || !iter._ba.Next() {
// 			return false
// 		}
// 		b, __a := iter._ba.Head()
// 		*iter.b = b
// 		iter.__a = toIterator(iter.ap, __a)
// 	}

// 	if iter.__a == nil || !iter.__a.Next() {
// 		iter.__a = nil
// 		return iter.Next()
// 	}

// 	return true
// }

// func (iter *iterator2[A, B]) FMap(f func(curie.IRI, curie.IRI, any) error) error {
// 	for iter.Next() {
// 		if err := f(iter.Head()); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
