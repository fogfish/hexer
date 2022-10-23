package hexagon

import (
	"math/rand"
	"time"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid"
	"github.com/fogfish/skiplist"
)

//
type Store struct {
	size   uint64
	random rand.Source
	spo    spo
	sop    sop
	pso    pso
	pos    pos
	osp    osp
	ops    ops
}

func New() *Store {
	rnd := rand.NewSource(time.Now().UnixNano())
	return &Store{
		random: rnd,
		spo:    newSPO(rnd),
		sop:    newSOP(rnd),
		pso:    newPSO(rnd),
		pos:    newPOS(rnd),
		osp:    newOSP(rnd),
		ops:    newOPS(rnd),
	}
}

func Put(store *Store, s, p curie.IRI, o any) {
	k := guid.L.K(guid.Clock)
	_po, _op := ensureForS(store, s)
	_so, _os := ensureForP(store, p)
	_sp, _ps := ensureForO(store, o)

	putO(store, _po, _so, s, p, o, k)
	putP(store, _op, _sp, s, p, o, k)
	putS(store, _os, _ps, s, p, o, k)
}

func ensureForS(store *Store, s curie.IRI) (_po, _op) {
	_po, has := skiplist.Lookup(store.spo, s)
	if !has {
		_po = newPO(store.random)
		skiplist.Put(store.spo, s, _po)
	}

	_op, has := skiplist.Lookup(store.sop, s)
	if !has {
		_op = newOP(store.random)
		skiplist.Put(store.sop, s, _op)
	}
	return _po, _op
}

func ensureForP(store *Store, p curie.IRI) (_so, _os) {
	_so, has := skiplist.Lookup(store.pso, p)
	if !has {
		_so = newSO(store.random)
		skiplist.Put(store.pso, p, _so)
	}

	_os, has := skiplist.Lookup(store.pos, p)
	if !has {
		_os = newOS(store.random)
		skiplist.Put(store.pos, p, _os)
	}
	return _so, _os
}

func ensureForO(store *Store, o any) (_sp, _ps) {
	_sp, has := skiplist.Lookup(store.osp, o)
	if !has {
		_sp = newSP(store.random)
		skiplist.Put(store.osp, o, _sp)
	}

	_ps, has := skiplist.Lookup(store.ops, o)
	if !has {
		_ps = newPS(store.random)
		skiplist.Put(store.ops, o, _ps)
	}
	return _sp, _ps
}

func putO(store *Store, _po _po, _so _so, s s, p p, o o, k k) {
	__o, has := skiplist.Lookup(_po, p)
	if !has {
		__o = newO(store.random)
		skiplist.Put(_po, p, __o)
		skiplist.Put(_so, s, __o)
	}

	skiplist.Put(__o, o, k)
}

func putP(store *Store, _op _op, _sp _sp, s s, p p, o o, k k) {
	__p, has := skiplist.Lookup(_sp, s)
	if !has {
		__p = newP(store.random)
		skiplist.Put(_op, o, __p)
		skiplist.Put(_sp, s, __p)
	}
	skiplist.Put(__p, p, k)
}

func putS(store *Store, _os _os, _ps _ps, s s, p p, o o, k k) {
	__s, has := skiplist.Lookup(_ps, p)
	if !has {
		__s = newS(store.random)
		skiplist.Put(_os, o, __s)
		skiplist.Put(_ps, p, __s)
	}
	skiplist.Put(__s, s, k)
}

func Query(store *Store, s *Predicate[s], p *Predicate[p], o *Predicate[o]) Iterator {
	q := pattern{store: store, s: s, p: p, o: o}
	_, iter := q.eval()
	return iter
}

//
type Node map[curie.IRI]any

func (node Node) Append(s, p curie.IRI, o any) error {
	if val, has := node[p]; !has {
		node[p] = o
	} else {
		switch v := val.(type) {
		case []any:
			node[p] = append(v, o)
		default:
			node[p] = []any{v, o}
		}
	}
	return nil
}

//
type Graph map[curie.IRI]Node

func (graph Graph) Append(s, p curie.IRI, o any) error {
	if val, has := graph[s]; !has {
		node := Node{}
		node.Append(s, p, o)
		graph[s] = node
	} else {
		val.Append(s, p, o)
		graph[s] = val
	}
	return nil
}
