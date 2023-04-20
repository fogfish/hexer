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

func (store *Store) sPO(ctx context.Context, q hexer.Pattern) hexer.StreamX {
	g := curie.IRI("a")

	key := spo{G: "sp|" + g, SP: encodeII(q.S.Value, "")}

	seq := NewIterator(store.spo, key)

	return &Unfold[spo]{seq: seq}

	// key := spo{G: "sp|a"}
	// seq := NewIterator(store.spo, key)
	// for seq.Next() {
	// 	fmt.Println(seq.Head())
	// }
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

// type Streamer struct {

// }

func Match(ctx context.Context, store *Store, q hexer.Pattern) (hexer.Strategy, hexer.StreamX) {
	strategy := q.Strategy()
	switch strategy {
	// x, o, _ ⇒ spo
	case 0510:
		return strategy, store.sPO(ctx, q)
	// // x, _, o ⇒ sop
	// case 0501:
	// 	return strategy, q.sOP()
	// // _, x, o ⇒ pos
	// case 0051:
	// 	return strategy, q.pOS()
	// // o, x, _ ⇒ pso
	// case 0150:
	// 	return strategy, q.pSO()
	// // o, _, x ⇒ osp
	// case 0105:
	// 	return strategy, q.oSP()
	// // _, o, x ⇒ ops
	// case 0015:
	// 	return strategy, q.oPS()

	// x, x, _ ⇒ spo
	case 0550:
		return strategy, store.sPO(ctx, q)
	// TODO: return strategy, q.spO()
	// // _, x, x ⇒ pos
	// case 0055:
	// 	return strategy, q.poS()
	// // x, _, x ⇒ sop
	// case 0505:
	// 	return strategy, q.soP()

	// x, _, _ ⇒ spo
	case 0500:
		return strategy, store.sPO(ctx, q)
		// TODO:	return strategy, q.sPO()
		// // _, x, _ ⇒ pso
		// case 0050:
		// 	return strategy, q.pSO()
		// // _, _, x ⇒ osp
		// case 0005:
		// 	return strategy, q.oSP()

		// // _, _, _ ⇒ spo
		// case 0000:
		// 	return strategy, q.spo()
	}

	return strategy, nil
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
