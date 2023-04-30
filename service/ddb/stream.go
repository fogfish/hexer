package ddb

import (
	"context"
	"fmt"

	"github.com/fogfish/curie"
	"github.com/fogfish/hexer"
)

type notSupported struct{ hexer.Pattern }

func (err notSupported) Error() string { return fmt.Sprintf("not supported %s", err.Pattern.Dump()) }
func (notSupported) NotSupported()     {}

func (store *Store) streamSPO(ctx context.Context, q hexer.Pattern) (hexer.Stream, error) {
	g := curie.IRI("a")
	key := spo{G: "sp|" + g}

	switch {
	case q.HintForS == hexer.HINT_MATCH && q.HintForP == hexer.HINT_NONE:
		key.SP = encodeII(q.S.Value, "")
	case q.HintForS == hexer.HINT_MATCH && q.HintForP == hexer.HINT_MATCH:
		key.SP = encodeII(q.S.Value, q.P.Value)
	case q.HintForS == hexer.HINT_MATCH && q.HintForP == hexer.HINT_FILTER_PREFIX:
		key.SP = encodeII(q.S.Value, q.P.Value)
	case q.HintForS == hexer.HINT_FILTER_PREFIX && q.HintForP == hexer.HINT_NONE:
		key.SP = encodeI(q.S.Value)
	default:
		return nil, &notSupported{q}
	}

	var stream hexer.Stream = &Unfold[spo]{
		seq: NewIterator(store.spo, key),
	}

	if q.O != nil {
		stream = hexer.NewFilterO(q.HintForO, q.O, stream)
	}

	return stream, nil
}

func (store *Store) streamSOP(ctx context.Context, q hexer.Pattern) (hexer.Stream, error) {
	g := curie.IRI("a")
	key := sop{G: "so|" + g}

	switch {
	case q.HintForS == hexer.HINT_MATCH && q.HintForO == hexer.HINT_NONE:
		key.SO = encodeII(q.S.Value, "")
	case q.HintForS == hexer.HINT_MATCH && q.HintForO == hexer.HINT_MATCH:
		key.SO = encodeIV(q.S.Value, q.O.Value)
	case q.HintForS == hexer.HINT_MATCH && q.HintForO == hexer.HINT_FILTER_PREFIX:
		key.SO = encodeIV(q.S.Value, q.O.Value)
	case q.HintForS == hexer.HINT_FILTER_PREFIX && q.HintForO == hexer.HINT_NONE:
		key.SO = encodeI(q.S.Value)
	default:
		return nil, &notSupported{q}
	}

	var stream hexer.Stream = &Unfold[sop]{
		seq: NewIterator(store.sop, key),
	}

	if q.P != nil {
		stream = hexer.NewFilterP(q.HintForP, q.P, stream)
	}

	return stream, nil
}

func (store *Store) streamPSO(ctx context.Context, q hexer.Pattern) (hexer.Stream, error) {
	g := curie.IRI("a")
	key := pso{G: "ps|" + g}

	switch {
	case q.HintForP == hexer.HINT_MATCH && q.HintForS == hexer.HINT_NONE:
		key.PS = encodeII(q.P.Value, "")
	case q.HintForP == hexer.HINT_MATCH && q.HintForS == hexer.HINT_MATCH:
		key.PS = encodeII(q.P.Value, q.S.Value)
	case q.HintForP == hexer.HINT_MATCH && q.HintForS == hexer.HINT_FILTER_PREFIX:
		key.PS = encodeII(q.P.Value, q.S.Value)
	case q.HintForP == hexer.HINT_FILTER_PREFIX && q.HintForS == hexer.HINT_NONE:
		key.PS = encodeI(q.P.Value)
	default:
		return nil, &notSupported{q}
	}

	var stream hexer.Stream = &Unfold[pso]{
		seq: NewIterator(store.pso, key),
	}

	if q.O != nil {
		stream = hexer.NewFilterO(q.HintForO, q.O, stream)
	}

	return stream, nil
}

func (store *Store) streamPOS(ctx context.Context, q hexer.Pattern) (hexer.Stream, error) {
	g := curie.IRI("a")
	key := pos{G: "po|" + g}

	switch {
	case q.HintForP == hexer.HINT_MATCH && q.HintForO == hexer.HINT_NONE:
		key.PO = encodeII(q.P.Value, "")
	case q.HintForP == hexer.HINT_MATCH && q.HintForO == hexer.HINT_MATCH:
		key.PO = encodeIV(q.P.Value, q.O.Value)
	case q.HintForP == hexer.HINT_MATCH && q.HintForO == hexer.HINT_FILTER_PREFIX:
		key.PO = encodeIV(q.P.Value, q.O.Value)
	case q.HintForP == hexer.HINT_FILTER_PREFIX && q.HintForO == hexer.HINT_NONE:
		key.PO = encodeI(q.P.Value)
	default:
		return nil, &notSupported{q}
	}

	var stream hexer.Stream = &Unfold[pos]{
		seq: NewIterator(store.pos, key),
	}

	if q.S != nil {
		stream = hexer.NewFilterS(q.HintForS, q.S, stream)
	}

	return stream, nil
}

func (store *Store) streamOSP(ctx context.Context, q hexer.Pattern) (hexer.Stream, error) {
	g := curie.IRI("a")
	key := osp{G: "os|" + g}

	switch {
	case q.HintForO == hexer.HINT_MATCH && q.HintForS == hexer.HINT_NONE:
		key.OS = encodeVI(q.O.Value, "")
	case q.HintForO == hexer.HINT_MATCH && q.HintForS == hexer.HINT_MATCH:
		key.OS = encodeVI(q.O.Value, q.S.Value)
	case q.HintForO == hexer.HINT_MATCH && q.HintForS == hexer.HINT_FILTER_PREFIX:
		key.OS = encodeVI(q.O.Value, q.S.Value)
	case q.HintForO == hexer.HINT_FILTER_PREFIX && q.HintForS == hexer.HINT_NONE:
		key.OS = encodeValue(q.O.Value)
	default:
		return nil, &notSupported{q}
	}

	var stream hexer.Stream = &Unfold[osp]{
		seq: NewIterator(store.osp, key),
	}

	if q.P != nil {
		stream = hexer.NewFilterP(q.HintForP, q.P, stream)
	}

	return stream, nil
}

func (store *Store) streamOPS(ctx context.Context, q hexer.Pattern) (hexer.Stream, error) {
	g := curie.IRI("a")
	key := ops{G: "op|" + g}

	switch {
	case q.HintForO == hexer.HINT_MATCH && q.HintForP == hexer.HINT_NONE:
		key.OP = encodeVI(q.O.Value, "")
	case q.HintForO == hexer.HINT_MATCH && q.HintForP == hexer.HINT_MATCH:
		key.OP = encodeVI(q.O.Value, q.P.Value)
	case q.HintForO == hexer.HINT_MATCH && q.HintForP == hexer.HINT_FILTER_PREFIX:
		key.OP = encodeVI(q.O.Value, q.P.Value)
	case q.HintForO == hexer.HINT_FILTER_PREFIX && q.HintForP == hexer.HINT_NONE:
		key.OP = encodeValue(q.O.Value)
	default:
		return nil, &notSupported{q}
	}

	var stream hexer.Stream = &Unfold[ops]{
		seq: NewIterator(store.ops, key),
	}

	if q.S != nil {
		stream = hexer.NewFilterS(q.HintForS, q.S, stream)
	}

	return stream, nil
}
