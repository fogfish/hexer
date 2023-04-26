package ephemeral

import (
	"github.com/fogfish/hexer"
)

type ispo struct{ *Iterator[s, p, o] }

func (i *ispo) Head() hexer.SPOCK {
	s, p, o := i.Iterator.Head()
	return hexer.SPOCK{S: s, P: p, O: o}
}

func (i *ispo) Next() bool { return i.Iterator.Next() }

func (i *ispo) FMap(f func(hexer.SPOCK) error) error {
	for i.Next() {
		if err := f(i.Head()); err != nil {
			return err
		}
	}
	return nil
}

func (store *Store) streamSPO(q hexer.Query) hexer.Stream {
	iter := &Iterator[s, p, o]{
		abc: toIterator(q.Pattern.S, store.spo),
		pb:  q.Pattern.P,
		pc:  q.Pattern.O,
	}

	return &ispo{iter}

	// switch {
	// case q.HintForS == hexer.HINT_MATCH && q.HintForP == hexer.HINT_NONE && q.HintForO == hexer.HINT_NONE:
	// 	iter := &BiIterator[p, o]{s: q.Pattern.S.Value, ap: q.Pattern.O}

	// 	if _po, has := skiplist.Lookup(store.spo, q.Pattern.S.Value); has {
	// 		iter.a = &iter.o
	// 		iter.b = &iter.p
	// 		iter._ba = toIterator(q.Pattern.P, _po)

	// 		return iter
	// 	}

	// 	return iter

	// case q.HintForS == hexer.HINT_MATCH && q.HintForP == hexer.HINT_MATCH && q.HintForO == hexer.HINT_NONE:
	// 	iter := &Iterator[o]{s: q.Pattern.S.Value, p: q.Pattern.P.Value}

	// 	if _po, has := skiplist.Lookup(store.spo, q.Pattern.S.Value); has {
	// 		if __o, has := skiplist.Lookup(_po, q.Pattern.P.Value); has {
	// 			iter.a = &iter.o
	// 			iter.__a = toIterator(q.Pattern.O, __o)

	// 			return iter
	// 		}
	// 	}

	// 	return iter

	// 	key.SP = encodeII(q.Pattern.S.Value, q.Pattern.P.Value)
	// case q.HintForS == hexer.HINT_MATCH && q.HintForP == hexer.HINT_FILTER_PREFIX:
	// 	key.SP = encodeII(q.Pattern.S.Value, q.Pattern.P.Value)
	// case q.HintForS == hexer.HINT_FILTER_PREFIX:
	// 	key.SP = encodeI(q.Pattern.S.Value)
	// default:
	// 	panic("spo xxx")
	// }

	// var stream hexer.Stream = &Unfold[spo]{
	// 	seq: NewIterator(store.spo, key),
	// }

	// switch {
	// case q.HintForO == hexer.HINT_MATCH:
	// 	panic("spo o xxx")
	// case q.HintForO == hexer.HINT_FILTER_PREFIX:
	// 	panic("spo o xxx")
	// case q.HintForO == hexer.HINT_FILTER:
	// 	panic("spo o xxx")
	// }

	// return stream

	// it := Iterator{spo: skiplist.Values(store.spo)}

	// for it.Next() {
	// 	h := it.Head()
	// 	fmt.Printf("%v\n", h)
	// }

	// return nil
}
