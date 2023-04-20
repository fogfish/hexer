package ddb

import (
	"context"
	"fmt"

	"github.com/fogfish/dynamo/v2"
	"github.com/fogfish/dynamo/v2/service/ddb"
	"github.com/fogfish/hexer"
)

/*

todo:
[x] dynamo cursor
[x] iter of dynamo seq (abstract spo, sop, xxx types so that it is implemented once)
3. iter of spo (set)  (abstract spo, sop, xxx types so that it is implemented once)

for iter.Next() {
	if err := f(iter.Head()); err != nil {
		return err
	}
}

How to abstract "Entity"?
s -> map[p]o ?


*/

type none string

func (none) MatchOpt() {}

type Seq[T dynamo.Thing] interface {
	Head() T
	Next() bool
}

func NewIterator[T dynamo.Thing](store *ddb.Storage[T], query T) Seq[T] {
	return &Iterator[T]{
		store:  store,
		query:  query,
		cursor: none(""),
	}
}

type Iterator[T dynamo.Thing] struct {
	store  *ddb.Storage[T]
	query  T
	cursor dynamo.MatchOpt
	seq    []T
}

func (iter *Iterator[T]) Head() T {
	return iter.seq[0]
}

func (iter *Iterator[T]) Next() bool {
	if iter.seq != nil && len(iter.seq) > 1 {
		iter.seq = iter.seq[1:]
		return true
	}

	if iter.cursor == nil {
		return false
	}

	var err error
	iter.seq, iter.cursor, err = iter.store.Match(context.TODO(),
		iter.query, iter.cursor, dynamo.Limit(2),
	)
	if err != nil {
		return false
	}

	if len(iter.seq) == 0 {
		return false
	}

	return true
}

type Unfold[T dynamo.Thing] struct {
	seq Seq[T]
	bag []hexer.SPOCK
}

func (unfold *Unfold[T]) Head() hexer.SPOCK {
	return unfold.bag[0]
}

func (unfold *Unfold[T]) Next() bool {
	if unfold.bag != nil && len(unfold.bag) > 1 {
		unfold.bag = unfold.bag[1:]
		return true
	}

	if !unfold.seq.Next() {
		return false
	}

	switch vv := any(unfold.seq.Head()).(type) {
	case interface{ ToSPOCK() []hexer.SPOCK }:
		unfold.bag = vv.ToSPOCK()
	default:
		fmt.Printf("==> %T\n", vv)
	}

	return true
}

/*
func IT(store *Store) *Iterator3 {
	return &Iterator3{store: store, eos: false}
}

type Iterator3 struct {
	store  *Store
	head   spo
	tail   []spo
	cursor *spo
	eos    bool
}

func (iter *Iterator3) Head() {
	sp, o := iter.head.SP, iter.head.O[0]
	iter.head.O = iter.head.O[1:]

	fmt.Printf("==> %v %v\n", sp, o)
}

func (iter *Iterator3) Next() bool {
	g := curie.IRI("a")

	if len(iter.head.O) == 0 && (iter.tail == nil || len(iter.tail) == 0) {
		if iter.eos {
			return false
		}

		key := spo{G: "sp|" + g}
		if iter.cursor != nil {
			key, iter.cursor = *iter.cursor, nil
		}

		fmt.Println("=====>>>>")
		seq, _, err := iter.store.spo.Match(context.Background(), key, dynamo.Limit(2))
		if err != nil {
			return false
		}
		iter.head = seq[0]
		iter.tail = seq[1:]

		if len(seq) == 2 {
			iter.cursor = &seq[len(seq)-1]
			iter.eos = false
		} else {
			iter.eos = true
		}

		return true
	}

	if len(iter.head.O) == 0 {
		iter.head = iter.tail[0]
		iter.tail = iter.tail[1:]
	}

	return true
}
*/
