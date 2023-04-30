package ephemeral

import (
	"github.com/fogfish/hexer"
)

func (store *Store) streamSPO(q hexer.Pattern) (hexer.Stream, error) {
	return newIterator[s, p, o](querySPO(q), store.spo), nil
}

func (store *Store) streamSOP(q hexer.Pattern) (hexer.Stream, error) {
	return newIterator[s, o, p](querySOP(q), store.sop), nil
}

func (store *Store) streamPSO(q hexer.Pattern) (hexer.Stream, error) {
	return newIterator[p, s, o](queryPSO(q), store.pso), nil
}

func (store *Store) streamPOS(q hexer.Pattern) (hexer.Stream, error) {
	return newIterator[p, o, s](queryPOS(q), store.pos), nil
}

func (store *Store) streamOSP(q hexer.Pattern) (hexer.Stream, error) {
	return newIterator[o, s, p](queryOSP(q), store.osp), nil
}

func (store *Store) streamOPS(q hexer.Pattern) (hexer.Stream, error) {
	return newIterator[o, p, s](queryOPS(q), store.ops), nil
}
