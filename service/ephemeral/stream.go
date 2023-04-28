package ephemeral

import (
	"github.com/fogfish/hexer"
)

func (store *Store) streamSPO(q hexer.Pattern) hexer.Stream {
	return newIterator[s, p, o](querySPO(q), store.spo)
}

func (store *Store) streamSOP(q hexer.Pattern) hexer.Stream {
	return newIterator[s, o, p](querySOP(q), store.sop)
}

func (store *Store) streamPSO(q hexer.Pattern) hexer.Stream {
	return newIterator[p, s, o](queryPSO(q), store.pso)
}

func (store *Store) streamPOS(q hexer.Pattern) hexer.Stream {
	return newIterator[p, o, s](queryPOS(q), store.pos)
}

func (store *Store) streamOSP(q hexer.Pattern) hexer.Stream {
	return newIterator[o, s, p](queryOSP(q), store.osp)
}

func (store *Store) streamOPS(q hexer.Pattern) hexer.Stream {
	return newIterator[o, p, s](queryOPS(q), store.ops)
}
