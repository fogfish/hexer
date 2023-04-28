package ephemeral

import (
	"github.com/fogfish/hexer"
	"github.com/fogfish/skiplist"
)

type hSPO hexer.Query

func (q hSPO) L1(list *skiplist.SkipList[s, _po]) Seq[s, _po] {
	return overIRI[s, _po](q.Pattern.S)(list)
}

func (q hSPO) L2(list *skiplist.SkipList[p, __o]) Seq[p, __o] {
	return overIRI[p, __o](q.Pattern.P)(list)
}

func (q hSPO) L3(list *skiplist.SkipList[o, k]) Seq[o, k] {
	return overXSD[o, k](q.Pattern.O)(list)
}

type hSOP hexer.Query

func (q hSOP) L1(list *skiplist.SkipList[s, _op]) Seq[s, _op] {
	return overIRI[s, _op](q.Pattern.S)(list)
}

func (q hSOP) L2(list *skiplist.SkipList[o, __p]) Seq[o, __p] {
	return overXSD[o, __p](q.Pattern.O)(list)
}

func (q hSOP) L3(list *skiplist.SkipList[p, k]) Seq[p, k] {
	return overIRI[p, k](q.Pattern.P)(list)
}

type hPSO hexer.Query

func (q hPSO) L1(list *skiplist.SkipList[p, _so]) Seq[p, _so] {
	return overIRI[p, _so](q.Pattern.P)(list)
}

func (q hPSO) L2(list *skiplist.SkipList[s, __o]) Seq[s, __o] {
	return overIRI[s, __o](q.Pattern.S)(list)
}

func (q hPSO) L3(list *skiplist.SkipList[o, k]) Seq[o, k] {
	return overXSD[o, k](q.Pattern.O)(list)
}

type hPOS hexer.Query

func (q hPOS) L1(list *skiplist.SkipList[p, _os]) Seq[p, _os] {
	return overIRI[p, _os](q.Pattern.P)(list)
}

func (q hPOS) L2(list *skiplist.SkipList[o, __p]) Seq[o, __p] {
	return overXSD[o, __p](q.Pattern.O)(list)
}

func (q hPOS) L3(list *skiplist.SkipList[s, k]) Seq[s, k] {
	return overIRI[s, k](q.Pattern.S)(list)
}

type hOPS hexer.Query

func (q hOPS) L1(list *skiplist.SkipList[o, _ps]) Seq[o, _ps] {
	return overXSD[o, _ps](q.Pattern.O)(list)
}

func (q hOPS) L2(list *skiplist.SkipList[p, __s]) Seq[p, __s] {
	return overIRI[p, __s](q.Pattern.P)(list)
}

func (q hOPS) L3(list *skiplist.SkipList[s, k]) Seq[s, k] {
	return overIRI[s, k](q.Pattern.S)(list)
}

type hOSP hexer.Query

func (q hOSP) L1(list *skiplist.SkipList[o, _ps]) Seq[o, _ps] {
	return overXSD[o, _ps](q.Pattern.O)(list)
}

func (q hOSP) L2(list *skiplist.SkipList[s, __p]) Seq[s, __p] {
	return overIRI[s, __p](q.Pattern.S)(list)
}

func (q hOSP) L3(list *skiplist.SkipList[p, k]) Seq[p, k] {
	return overIRI[p, k](q.Pattern.P)(list)
}

func (store *Store) streamSPO(q hexer.Query) hexer.Stream {
	return NewIterator[s, p, o](
		hSPO(q),
		store.spo,
		// q.Pattern.S,
		// q.Pattern.P,
		// q.Pattern.O,
		func(s s, p p, o o) (hexer.SPOCK, bool) {
			return hexer.SPOCK{S: s, P: p, O: o}, true
		},
	)
}

func (store *Store) streamSOP(q hexer.Query) hexer.Stream {
	builder := func(s s, o o, p p) (hexer.SPOCK, bool) {
		return hexer.SPOCK{S: s, P: p, O: o}, true
	}

	if q.Pattern.O != nil {
		domain := q.Pattern.O.Value.XSDType()
		builder = func(s s, o o, p p) (hexer.SPOCK, bool) {
			return hexer.SPOCK{S: s, P: p, O: o}, o.XSDType() == domain
		}
	}

	return NewIterator[s, o, p](
		hSOP(q),
		store.sop,
		// q.Pattern.S,
		// q.Pattern.O,
		// q.Pattern.P,
		builder,
	)
}

func (store *Store) streamPSO(q hexer.Query) hexer.Stream {
	return NewIterator[p, s, o](
		hPSO(q),
		store.pso,
		// q.Pattern.P,
		// q.Pattern.S,
		// q.Pattern.O,
		func(p p, s s, o o) (hexer.SPOCK, bool) {
			return hexer.SPOCK{S: s, P: p, O: o}, true
		},
	)
}

func (store *Store) streamPOS(q hexer.Query) hexer.Stream {
	return NewIterator[p, o, s](
		hPOS(q),
		store.pos,
		// q.Pattern.P,
		// q.Pattern.O,
		// q.Pattern.S,
		func(p p, o o, s s) (hexer.SPOCK, bool) {
			return hexer.SPOCK{S: s, P: p, O: o}, true
		},
	)
}

func (store *Store) streamOSP(q hexer.Query) hexer.Stream {
	return NewIterator[o, s, p](
		hOSP(q),
		store.osp,
		// q.Pattern.O,
		// q.Pattern.S,
		// q.Pattern.P,
		func(o o, s s, p p) (hexer.SPOCK, bool) {
			return hexer.SPOCK{S: s, P: p, O: o}, true
		},
	)
}

func (store *Store) streamOPS(q hexer.Query) hexer.Stream {
	return NewIterator[o, p, s](
		hOPS(q),
		store.ops,
		// q.Pattern.O,
		// q.Pattern.P,
		// q.Pattern.S,
		func(o o, p p, s s) (hexer.SPOCK, bool) {
			return hexer.SPOCK{S: s, P: p, O: o}, true
		},
	)
}
