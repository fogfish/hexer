package ephemeral

import (
	"github.com/fogfish/hexer"
	"github.com/fogfish/skiplist"
)

type Helper[A, B, C any] interface {
	L1(*skiplist.SkipList[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]) Seq[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]
	L2(*skiplist.SkipList[B, *skiplist.SkipList[C, k]]) Seq[B, *skiplist.SkipList[C, k]]
	L3(*skiplist.SkipList[C, k]) Seq[C, k]
	ToSPOCK(A, B, C) hexer.SPOCK
}

type Iterator[A, B, C any] struct {
	a   A
	b   B
	c   C
	abc Seq[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]]
	_bc Seq[B, *skiplist.SkipList[C, k]]
	__c Seq[C, k]
	hlp Helper[A, B, C]
}

func NewIterator[A, B, C any](
	hlp Helper[A, B, C],
	seq *skiplist.SkipList[A, *skiplist.SkipList[B, *skiplist.SkipList[C, k]]],
) *Iterator[A, B, C] {
	return &Iterator[A, B, C]{
		hlp: hlp,
		abc: hlp.L1(seq),
	}
}

func (iter *Iterator[A, B, C]) Head() hexer.SPOCK {
	return iter.hlp.ToSPOCK(iter.a, iter.b, iter.c)
}

func (iter *Iterator[A, B, C]) Next() bool {
	if iter._bc == nil {
		if iter.abc == nil || !iter.abc.Next() {
			return false
		}
		a, _bc := iter.abc.Head()
		iter.a = a
		iter._bc = iter.hlp.L2(_bc)
	}

	if iter.__c == nil {
		if iter._bc == nil || !iter._bc.Next() {
			iter._bc = nil
			return iter.Next()
		}

		b, __c := iter._bc.Head()
		iter.b = b
		iter.__c = iter.hlp.L3(__c)
	}

	if iter.__c == nil || !iter.__c.Next() {
		iter.__c = nil
		return iter.Next()
	}

	iter.c, _ = iter.__c.Head()

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
