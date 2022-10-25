//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/hexagon
//

package hexagon

import (
	"math/rand"
	"time"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid"
	"github.com/fogfish/skiplist"
)

// Store is the instance of knowledge storage
type Store struct {
	size   int
	random rand.Source
	spo    spo
	sop    sop
	pso    pso
	pos    pos
	osp    osp
	ops    ops
}

// Create new instance of knowledge storage
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

// Size returns number of knowledge statements in the store
func Size(store *Store) int {
	return store.size
}

// Put new knowledge statement into the store
func Put[T DataType](store *Store, s, p curie.IRI, o T) {
	k := guid.L.K(guid.Clock)
	_po, _op := ensureForS(store, s)
	_so, _os := ensureForP(store, p)
	_sp, _ps := ensureForO(store, o)

	putO(store, _po, _so, s, p, o, k)
	putP(store, _op, _sp, s, p, o, k)
	putS(store, _os, _ps, s, p, o, k)
	store.size++
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

// Match knowledge statements to the pattern and return stream of knowledge statements
func Match(store *Store, s *Predicate[s], p *Predicate[p], o *Predicate[o]) Stream {
	q := pattern{store: store, s: s, p: p, o: o}
	_, iter := q.eval()
	return iter
}
