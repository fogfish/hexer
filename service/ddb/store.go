package ddb

import (
	"context"

	"github.com/fogfish/curie"
	"github.com/fogfish/dynamo/v2"
	"github.com/fogfish/dynamo/v2/service/ddb"
	"github.com/fogfish/hexer"
)

type Store struct {
	spo *ddb.Storage[spo]
	sop *ddb.Storage[sop]
	pso *ddb.Storage[pso]
	pos *ddb.Storage[pos]
	osp *ddb.Storage[osp]
	ops *ddb.Storage[ops]
}

func New(connector string, opts ...dynamo.Option) (*Store, error) {
	spo, err := ddb.New[spo](connector, opts...)
	if err != nil {
		return nil, err
	}

	sop, err := ddb.New[sop](connector, opts...)
	if err != nil {
		return nil, err
	}

	pso, err := ddb.New[pso](connector, opts...)
	if err != nil {
		return nil, err
	}

	pos, err := ddb.New[pos](connector, opts...)
	if err != nil {
		return nil, err
	}

	osp, err := ddb.New[osp](connector, opts...)
	if err != nil {
		return nil, err
	}

	ops, err := ddb.New[ops](connector, opts...)
	if err != nil {
		return nil, err
	}

	return &Store{
		spo: spo,
		sop: sop,
		pso: pso,
		pos: pos,
		osp: osp,
		ops: ops,
	}, nil
}

func Add(ctx context.Context, store *Store, bag hexer.Bag) (hexer.Bag, error) {
	for i, spock := range bag {
		if err := Put(ctx, store, spock); err != nil {
			return bag[i:], err
		}
	}

	return nil, nil
}

func Put(ctx context.Context, store *Store, spock hexer.SPOCK) error {
	g := curie.IRI("a")

	seq := []Writer{
		encodeSPO(g, spock),
		encodeSOP(g, spock),
		encodePOS(g, spock),
		encodePSO(g, spock),
		encodeOPS(g, spock),
		encodeOSP(g, spock),
	}

	for i := 0; i < len(seq); i++ {
		if err := seq[i].Put(ctx, store); err != nil {
			for k := 0; k < i; k++ {
				if err := seq[k].Cut(ctx, store); err != nil {
					// TODO: log error
				}
			}
			return err
		}
	}

	return nil
}

func Match(ctx context.Context, store *Store, q hexer.Pattern) hexer.Stream {
	switch q.Strategy {
	case hexer.STRATEGY_SPO:
		return store.streamSPO(ctx, q)
	case hexer.STRATEGY_SOP:
		return store.streamSOP(ctx, q)
	case hexer.STRATEGY_PSO:
		return store.streamPSO(ctx, q)
	case hexer.STRATEGY_POS:
		return store.streamPOS(ctx, q)
	case hexer.STRATEGY_OSP:
		return store.streamOSP(ctx, q)
	case hexer.STRATEGY_OPS:
		return store.streamOPS(ctx, q)
	default:
		panic("xxx")
	}
}

// func Get(ctx context.Context, store *Store, spock hexer.SPOCK[string]) SPO {
// 	key := SPO{
// 		G:  "g:a",
// 		SP: string(spock.S) + "#" + string(spock.P),
// 	}

// 	spo, err := store.spo.Get(ctx, key)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return spo
// }
