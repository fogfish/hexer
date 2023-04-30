package ephemeral

import (
	"math/rand"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid/v2"
	"github.com/fogfish/hexer/internal/ord"
	"github.com/fogfish/hexer/xsd"
	"github.com/fogfish/skiplist"
)

// components of <s,p,o,c,k> triple
type s = curie.IRI // subject
type p = curie.IRI // predicate
type o = xsd.Value // object
type c = float64   // credibility
type k = guid.K    // k-order

// index types for 3rd faction
type __s = *skiplist.SkipList[s, k]
type __p = *skiplist.SkipList[p, k]
type __o = *skiplist.SkipList[o, k]

// index types for 2nd faction
type _po = *skiplist.SkipList[p, __o]
type _op = *skiplist.SkipList[o, __p]
type _so = *skiplist.SkipList[s, __o]
type _os = *skiplist.SkipList[o, __s]
type _sp = *skiplist.SkipList[s, __p]
type _ps = *skiplist.SkipList[p, __s]

// triple indexes
type spo = *skiplist.SkipList[s, _po]
type sop = *skiplist.SkipList[s, _op]
type pso = *skiplist.SkipList[p, _so]
type pos = *skiplist.SkipList[p, _os]
type osp = *skiplist.SkipList[o, _sp]
type ops = *skiplist.SkipList[o, _ps]

// allocators for indexes
func newS(rnd rand.Source) __s { return skiplist.New[s, k](ord.IRI, rnd) }
func newP(rnd rand.Source) __p { return skiplist.New[p, k](ord.IRI, rnd) }
func newO(rnd rand.Source) __o { return skiplist.New[o, k](ord.XSD, rnd) }

func newPO(rnd rand.Source) _po { return skiplist.New[p, __o](ord.IRI, rnd) }
func newOP(rnd rand.Source) _op { return skiplist.New[o, __p](ord.XSD, rnd) }
func newSO(rnd rand.Source) _so { return skiplist.New[s, __o](ord.IRI, rnd) }
func newOS(rnd rand.Source) _os { return skiplist.New[o, __s](ord.XSD, rnd) }
func newSP(rnd rand.Source) _sp { return skiplist.New[s, __p](ord.IRI, rnd) }
func newPS(rnd rand.Source) _ps { return skiplist.New[p, __s](ord.IRI, rnd) }

func newSPO(rnd rand.Source) spo { return skiplist.New[s, _po](ord.IRI, rnd) }
func newSOP(rnd rand.Source) sop { return skiplist.New[s, _op](ord.IRI, rnd) }
func newPSO(rnd rand.Source) pso { return skiplist.New[p, _so](ord.IRI, rnd) }
func newPOS(rnd rand.Source) pos { return skiplist.New[p, _os](ord.IRI, rnd) }
func newOSP(rnd rand.Source) osp { return skiplist.New[o, _sp](ord.XSD, rnd) }
func newOPS(rnd rand.Source) ops { return skiplist.New[o, _ps](ord.XSD, rnd) }
