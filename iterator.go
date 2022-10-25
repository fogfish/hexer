//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/hexagon
//

package hexagon

import (
	"github.com/fogfish/curie"
	"github.com/fogfish/skiplist"
)

/*

iterator build stream of <s,p,o> iterating over single dimension

the pointer `a` holds address of the field s, p, o, depending on the config
*/
type iterator[A any] struct {
	s   s
	p   p
	o   o
	a   *A
	__a *skiplist.Iterator[A, k]
}

func (iter *iterator[A]) Head() (s, p, o) {
	*iter.a, _ = iter.__a.Head()
	return iter.s, iter.p, iter.o
}

func (iter *iterator[A]) Next() bool {
	return iter.__a != nil && iter.__a.Next()
}

func (iter *iterator[T]) FMap(f func(curie.IRI, curie.IRI, any) error) error {
	for iter.Next() {
		if err := f(iter.Head()); err != nil {
			return err
		}
	}
	return nil
}

/*

iterator2 build stream of <s,p,o> iterating over two dimension
*/
type iterator2[B, A any] struct {
	s   s
	p   p
	o   o
	a   *A
	b   *B
	ap  *Predicate[A]
	_ba *skiplist.Iterator[B, *skiplist.SkipList[A, k]]
	__a *skiplist.Iterator[A, k]
}

func (iter *iterator2[A, B]) Head() (s, p, o) {
	*iter.a, _ = iter.__a.Head()
	return iter.s, iter.p, iter.o
}

func (iter *iterator2[A, B]) Next() bool {
	if iter.__a == nil {
		if iter._ba == nil || !iter._ba.Next() {
			return false
		}
		b, __a := iter._ba.Head()
		*iter.b = b
		iter.__a = toIterator(iter.ap, __a)
	}

	if iter.__a == nil || !iter.__a.Next() {
		iter.__a = nil
		return iter.Next()
	}

	return true
}

func (iter *iterator2[A, B]) FMap(f func(curie.IRI, curie.IRI, any) error) error {
	for iter.Next() {
		if err := f(iter.Head()); err != nil {
			return err
		}
	}
	return nil
}

/*

iterator3 build stream of <s,p,o> iterating over three dimension
*/
type iterator3 struct {
	s   s
	p   p
	o   o
	spo *skiplist.Iterator[s, _po]
	_po *skiplist.Iterator[p, __o]
	__o *skiplist.Iterator[o, k]
}

func (iter *iterator3) Head() (s, p, o) {
	iter.o, _ = iter.__o.Head()
	return iter.s, iter.p, iter.o
}

func (iter *iterator3) Next() bool {
	if iter._po == nil {
		if iter.spo == nil || !iter.spo.Next() {
			return false
		}
		s, _po := iter.spo.Head()
		iter.s = s
		iter._po = skiplist.Values(_po)
	}

	if iter.__o == nil {
		if iter._po == nil || !iter._po.Next() {
			iter._po = nil
			return iter.Next()
		}

		p, __o := iter._po.Head()
		iter.p = p
		iter.__o = skiplist.Values(__o)
	}

	if iter.__o == nil || !iter.__o.Next() {
		iter.__o = nil
		return iter.Next()
	}

	return true
}

func (iter *iterator3) FMap(f func(curie.IRI, curie.IRI, any) error) error {
	for iter.Next() {
		if err := f(iter.Head()); err != nil {
			return err
		}
	}
	return nil
}
