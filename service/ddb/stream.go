package ddb

import (
	"context"
	"strings"

	"github.com/fogfish/curie"
	"github.com/fogfish/hexer"
)

func (store *Store) streamSPO(ctx context.Context, q hexer.Query) hexer.Stream {
	g := curie.IRI("a")
	key := spo{G: "sp|" + g}

	switch {
	case q.HintForS == hexer.HINT_MATCH && q.HintForP == hexer.HINT_NONE:
		key.SP = encodeII(q.Pattern.S.Value, "")
	case q.HintForS == hexer.HINT_MATCH && q.HintForP == hexer.HINT_MATCH:
		key.SP = encodeII(q.Pattern.S.Value, q.Pattern.P.Value)
	case q.HintForS == hexer.HINT_MATCH && q.HintForP == hexer.HINT_FILTER_PREFIX:
		key.SP = encodeII(q.Pattern.S.Value, q.Pattern.P.Value)
	case q.HintForS == hexer.HINT_FILTER_PREFIX:
		key.SP = encodeI(q.Pattern.S.Value)
	default:
		panic("spo xxx")
	}

	var stream hexer.Stream = &Unfold[spo]{
		seq: NewIterator(store.spo, key),
	}

	switch {
	case q.HintForO == hexer.HINT_MATCH:
		panic("spo o xxx")
	case q.HintForO == hexer.HINT_FILTER_PREFIX:
		panic("spo o xxx")
	case q.HintForO == hexer.HINT_FILTER:
		panic("spo o xxx")
	}

	return stream
}

func (store *Store) streamSOP(ctx context.Context, q hexer.Query) hexer.Stream {
	g := curie.IRI("a")
	key := sop{G: "so|" + g}

	switch {
	case q.HintForS == hexer.HINT_MATCH && q.HintForO == hexer.HINT_NONE:
		key.SO = encodeII(q.Pattern.S.Value, "")
	case q.HintForS == hexer.HINT_MATCH && q.HintForO == hexer.HINT_MATCH:
		key.SO = encodeIV(q.Pattern.S.Value, q.Pattern.O.Value)
	case q.HintForS == hexer.HINT_MATCH && q.HintForO == hexer.HINT_FILTER_PREFIX:
		key.SO = encodeIV(q.Pattern.S.Value, q.Pattern.O.Value)
	case q.HintForS == hexer.HINT_FILTER_PREFIX:
		key.SO = encodeI(q.Pattern.S.Value)
	default:
		panic("sop so xxx")
	}

	var stream hexer.Stream = &Unfold[sop]{
		seq: NewIterator(store.sop, key),
	}

	switch {
	case q.HintForP == hexer.HINT_MATCH:
		stream = hexer.NewFilter(
			func(spock hexer.SPOCK) bool { return spock.P == q.Pattern.P.Value },
			stream,
		)
	case q.HintForP == hexer.HINT_FILTER_PREFIX:
		stream = hexer.NewFilter(
			func(spock hexer.SPOCK) bool { return strings.HasPrefix(string(spock.P), string(q.Pattern.P.Value)) },
			stream,
		)
	case q.HintForP == hexer.HINT_FILTER:
		panic("sop p xxx")
	}

	return stream
}

func (store *Store) streamPSO(ctx context.Context, q hexer.Query) hexer.Stream {
	g := curie.IRI("a")
	key := pso{G: "ps|" + g}

	switch {
	case q.HintForP == hexer.HINT_MATCH && q.HintForS == hexer.HINT_NONE:
		key.PS = encodeII(q.Pattern.P.Value, "")
	case q.HintForP == hexer.HINT_MATCH && q.HintForS == hexer.HINT_MATCH:
		key.PS = encodeII(q.Pattern.P.Value, q.Pattern.S.Value)
	case q.HintForP == hexer.HINT_MATCH && q.HintForS == hexer.HINT_FILTER_PREFIX:
		key.PS = encodeII(q.Pattern.P.Value, q.Pattern.S.Value)
	case q.HintForP == hexer.HINT_FILTER_PREFIX:
		key.PS = encodeI(q.Pattern.P.Value)
	default:
		panic("pso xxx")
	}

	var stream hexer.Stream = &Unfold[pso]{
		seq: NewIterator(store.pso, key),
	}

	switch {
	case q.HintForO == hexer.HINT_MATCH:
		panic("spo o xxx")
	case q.HintForO == hexer.HINT_FILTER_PREFIX:
		panic("spo o xxx")
	case q.HintForO == hexer.HINT_FILTER:
		panic("spo o xxx")
	}

	return stream
}

func (store *Store) streamPOS(ctx context.Context, q hexer.Query) hexer.Stream {
	g := curie.IRI("a")
	key := pos{G: "po|" + g}

	switch {
	case q.HintForP == hexer.HINT_MATCH && q.HintForO == hexer.HINT_NONE:
		key.PO = encodeII(q.Pattern.S.Value, "")
	case q.HintForP == hexer.HINT_MATCH && q.HintForO == hexer.HINT_MATCH:
		key.PO = encodeIV(q.Pattern.S.Value, q.Pattern.O.Value)
	case q.HintForP == hexer.HINT_MATCH && q.HintForO == hexer.HINT_FILTER_PREFIX:
		key.PO = encodeIV(q.Pattern.S.Value, q.Pattern.O.Value)
	case q.HintForP == hexer.HINT_FILTER_PREFIX:
		key.PO = encodeI(q.Pattern.S.Value)
	default:
		panic("pos so xxx")
	}

	var stream hexer.Stream = &Unfold[pos]{
		seq: NewIterator(store.pos, key),
	}

	switch {
	case q.HintForS == hexer.HINT_MATCH:
		stream = hexer.NewFilter(
			func(spock hexer.SPOCK) bool { return spock.S == q.Pattern.S.Value },
			stream,
		)
	case q.HintForS == hexer.HINT_FILTER_PREFIX:
		stream = hexer.NewFilter(
			func(spock hexer.SPOCK) bool { return strings.HasPrefix(string(spock.S), string(q.Pattern.S.Value)) },
			stream,
		)
	case q.HintForS == hexer.HINT_FILTER:
		panic("pos s xxx")
	}

	return stream
}

func (store *Store) streamOSP(ctx context.Context, q hexer.Query) hexer.Stream {
	g := curie.IRI("a")
	key := osp{G: "os|" + g}

	switch {
	case q.HintForO == hexer.HINT_MATCH && q.HintForS == hexer.HINT_NONE:
		key.OS = encodeVI(q.Pattern.O.Value, "")
	case q.HintForO == hexer.HINT_MATCH && q.HintForS == hexer.HINT_MATCH:
		key.OS = encodeVI(q.Pattern.O.Value, q.Pattern.S.Value)
	case q.HintForO == hexer.HINT_MATCH && q.HintForS == hexer.HINT_FILTER_PREFIX:
		key.OS = encodeVI(q.Pattern.O.Value, q.Pattern.S.Value)
	case q.HintForO == hexer.HINT_FILTER_PREFIX:
		key.OS = encodeValue(q.Pattern.O.Value)
	default:
		panic("osp os xxx")
	}

	var stream hexer.Stream = &Unfold[osp]{
		seq: NewIterator(store.osp, key),
	}

	switch {
	case q.HintForP == hexer.HINT_MATCH:
		stream = hexer.NewFilter(
			func(spock hexer.SPOCK) bool { return spock.P == q.Pattern.P.Value },
			stream,
		)
	case q.HintForP == hexer.HINT_FILTER_PREFIX:
		stream = hexer.NewFilter(
			func(spock hexer.SPOCK) bool { return strings.HasPrefix(string(spock.P), string(q.Pattern.P.Value)) },
			stream,
		)
	case q.HintForP == hexer.HINT_FILTER:
		panic("sop p xxx")
	}

	return stream
}

func (store *Store) streamOPS(ctx context.Context, q hexer.Query) hexer.Stream {
	g := curie.IRI("a")
	key := ops{G: "op|" + g}

	switch {
	case q.HintForO == hexer.HINT_MATCH && q.HintForP == hexer.HINT_NONE:
		key.OP = encodeVI(q.Pattern.O.Value, "")
	case q.HintForO == hexer.HINT_MATCH && q.HintForP == hexer.HINT_MATCH:
		key.OP = encodeVI(q.Pattern.O.Value, q.Pattern.P.Value)
	case q.HintForO == hexer.HINT_MATCH && q.HintForP == hexer.HINT_FILTER_PREFIX:
		key.OP = encodeVI(q.Pattern.O.Value, q.Pattern.S.Value)
	case q.HintForO == hexer.HINT_FILTER_PREFIX:
		key.OP = encodeValue(q.Pattern.O.Value)
	default:
		panic("ops op xxx")
	}

	var stream hexer.Stream = &Unfold[ops]{
		seq: NewIterator(store.ops, key),
	}

	switch {
	case q.HintForS == hexer.HINT_MATCH:
		stream = hexer.NewFilter(
			func(spock hexer.SPOCK) bool { return spock.S == q.Pattern.S.Value },
			stream,
		)
	case q.HintForS == hexer.HINT_FILTER_PREFIX:
		stream = hexer.NewFilter(
			func(spock hexer.SPOCK) bool { return strings.HasPrefix(string(spock.S), string(q.Pattern.S.Value)) },
			stream,
		)
	case q.HintForS == hexer.HINT_FILTER:
		panic("ops s xxx")
	}

	return stream
}
