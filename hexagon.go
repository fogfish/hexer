package hexagon

import (
	"github.com/fogfish/curie"
	"github.com/fogfish/guid"
	"github.com/fogfish/skiplist"
)

//
type Store struct {
	size uint64
	spo  spo
	sop  sop
	pso  pso
	pos  pos
	osp  osp
	ops  ops
}

func New() *Store {
	return &Store{
		spo: newSPO(),
		sop: newSOP(),
		pso: newPSO(),
		pos: newPOS(),
		osp: newOSP(),
		ops: newOPS(),
	}
}

func Put(store *Store, s, p curie.IRI, o any) {
	k := guid.L.K(guid.Clock)
	_po, _op := ensureForS(store, s)
	_so, _os := ensureForP(store, p)
	_sp, _ps := ensureForO(store, o)

	putO(_po, _so, s, p, o, k)
	putP(_op, _sp, s, p, o, k)
	putS(_os, _ps, s, p, o, k)
}

func ensureForS(store *Store, s curie.IRI) (_po, _op) {
	_po, has := skiplist.Lookup(store.spo, s)
	if !has {
		_po = newPO()
		skiplist.Put(store.spo, s, _po)
	}

	_op, has := skiplist.Lookup(store.sop, s)
	if !has {
		_op = newOP()
		skiplist.Put(store.sop, s, _op)
	}
	return _po, _op
}

func ensureForP(store *Store, p curie.IRI) (_so, _os) {
	_so, has := skiplist.Lookup(store.pso, p)
	if !has {
		_so = newSO()
		skiplist.Put(store.pso, p, _so)
	}

	_os, has := skiplist.Lookup(store.pos, p)
	if !has {
		_os = newOS()
		skiplist.Put(store.pos, p, _os)
	}
	return _so, _os
}

func ensureForO(store *Store, o any) (_sp, _ps) {
	_sp, has := skiplist.Lookup(store.osp, o)
	if !has {
		_sp = newSP()
		skiplist.Put(store.osp, o, _sp)
	}

	_ps, has := skiplist.Lookup(store.ops, o)
	if !has {
		_ps = newPS()
		skiplist.Put(store.ops, o, _ps)
	}
	return _sp, _ps
}

func putO(_po _po, _so _so, s s, p p, o o, k k) {
	__o, has := skiplist.Lookup(_po, p)
	if !has {
		__o = newO()
		skiplist.Put(_po, p, __o)
		skiplist.Put(_so, s, __o)
	}

	skiplist.Put(__o, o, k)
}

func putP(_op _op, _sp _sp, s s, p p, o o, k k) {
	__p, has := skiplist.Lookup(_sp, s)
	if !has {
		__p = newP()
		skiplist.Put(_op, o, __p)
		skiplist.Put(_sp, s, __p)
	}
	skiplist.Put(__p, p, k)
}

func putS(_os _os, _ps _ps, s s, p p, o o, k k) {
	__s, has := skiplist.Lookup(_ps, p)
	if !has {
		__s = newS()
		skiplist.Put(_os, o, __s)
		skiplist.Put(_ps, p, __s)
	}
	skiplist.Put(__s, s, k)
}

func Query(store *Store, s *Predicate[s], p *Predicate[s], o *Predicate[o]) Iterator {
	q := pattern{store: store, s: s, p: p, o: o}
	_, iter := q.eval()
	return iter
}
