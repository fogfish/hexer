package ephemeral

import (
	"github.com/fogfish/hexer"
)

func (store *Store) streamSPO(q hexer.Query) hexer.Stream {
	return NewIterator[s, p, o](
		querySPO(q),
		store.spo,
		// q.Pattern.S,
		// q.Pattern.P,
		// q.Pattern.O,
		// func(s s, p p, o o) (hexer.SPOCK, bool) {
		// 	return hexer.SPOCK{S: s, P: p, O: o}, true
		// },
	)
}

func (store *Store) streamSOP(q hexer.Query) hexer.Stream {
	// builder := func(s s, o o, p p) (hexer.SPOCK, bool) {
	// 	return hexer.SPOCK{S: s, P: p, O: o}, true
	// }

	// if q.Pattern.O != nil {
	// 	domain := q.Pattern.O.Value.XSDType()
	// 	builder = func(s s, o o, p p) (hexer.SPOCK, bool) {
	// 		return hexer.SPOCK{S: s, P: p, O: o}, o.XSDType() == domain
	// 	}
	// }

	return NewIterator[s, o, p](
		querySOP(q),
		store.sop,
		// q.Pattern.S,
		// q.Pattern.O,
		// q.Pattern.P,
		// builder,
	)
}

func (store *Store) streamPSO(q hexer.Query) hexer.Stream {
	return NewIterator[p, s, o](
		queryPSO(q),
		store.pso,
		// q.Pattern.P,
		// q.Pattern.S,
		// q.Pattern.O,
		// func(p p, s s, o o) (hexer.SPOCK, bool) {
		// 	return hexer.SPOCK{S: s, P: p, O: o}, true
		// },
	)
}

func (store *Store) streamPOS(q hexer.Query) hexer.Stream {
	return NewIterator[p, o, s](
		queryPOS(q),
		store.pos,
		// q.Pattern.P,
		// q.Pattern.O,
		// q.Pattern.S,
		// func(p p, o o, s s) (hexer.SPOCK, bool) {
		// 	return hexer.SPOCK{S: s, P: p, O: o}, true
		// },
	)
}

func (store *Store) streamOSP(q hexer.Query) hexer.Stream {
	return NewIterator[o, s, p](
		queryOSP(q),
		store.osp,
		// q.Pattern.O,
		// q.Pattern.S,
		// q.Pattern.P,
		// func(o o, s s, p p) (hexer.SPOCK, bool) {
		// 	return hexer.SPOCK{S: s, P: p, O: o}, true
		// },
	)
}

func (store *Store) streamOPS(q hexer.Query) hexer.Stream {
	return NewIterator[o, p, s](
		queryOPS(q),
		store.ops,
		// q.Pattern.O,
		// q.Pattern.P,
		// q.Pattern.S,
		// func(o o, p p, s s) (hexer.SPOCK, bool) {
		// 	return hexer.SPOCK{S: s, P: p, O: o}, true
		// },
	)
}
