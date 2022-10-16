package hexagon

import (
	"github.com/fogfish/curie"
	"github.com/fogfish/guid"
	"github.com/fogfish/skiplist"
)

type s = curie.IRI
type p = curie.IRI
type o = any
type k = guid.LID

type __s = *skiplist.SkipList[s, k]
type __p = *skiplist.SkipList[p, k]
type __o = *skiplist.SkipList[o, k]

type _po = *skiplist.SkipList[p, __o]
type _op = *skiplist.SkipList[o, __p]
type _so = *skiplist.SkipList[s, __o]
type _os = *skiplist.SkipList[o, __s]
type _sp = *skiplist.SkipList[s, __p]
type _ps = *skiplist.SkipList[p, __s]

type spo = *skiplist.SkipList[s, _po]
type sop = *skiplist.SkipList[s, _op]
type pso = *skiplist.SkipList[p, _so]
type pos = *skiplist.SkipList[p, _os]
type osp = *skiplist.SkipList[o, _sp]
type ops = *skiplist.SkipList[o, _ps]

func newS() __s { return skiplist.New[s, k](ordIRI) }
func newP() __p { return skiplist.New[p, k](ordIRI) }
func newO() __o { return skiplist.New[o, k](ordAny) }

func newPO() _po { return skiplist.New[p, __o](ordIRI) }
func newOP() _op { return skiplist.New[o, __p](ordAny) }
func newSO() _so { return skiplist.New[s, __o](ordIRI) }
func newOS() _os { return skiplist.New[o, __s](ordAny) }
func newSP() _sp { return skiplist.New[s, __p](ordIRI) }
func newPS() _ps { return skiplist.New[p, __s](ordIRI) }

func newSPO() spo { return skiplist.New[s, _po](ordIRI) }
func newSOP() sop { return skiplist.New[s, _op](ordIRI) }
func newPSO() pso { return skiplist.New[p, _so](ordIRI) }
func newPOS() pos { return skiplist.New[p, _os](ordIRI) }
func newOSP() osp { return skiplist.New[o, _sp](ordAny) }
func newOPS() ops { return skiplist.New[o, _ps](ordAny) }
