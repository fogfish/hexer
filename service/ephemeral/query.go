package ephemeral

import (
	"strings"

	"github.com/fogfish/curie"
	"github.com/fogfish/hexer"
	"github.com/fogfish/hexer/xsd"
	"github.com/fogfish/skiplist"
)

// Each query results with sequence of "elements".
// This interface defines generic sequence, abstracting skiplist.Iterator
type Seq[K, V any] interface {
	Head() (K, V)
	Next() bool
}

// helper function to query the skiplist where key is curie.IRI
func queryIRI[A, B any](
	pred *hexer.Predicate[curie.IRI],
	list *skiplist.SkipList[curie.IRI, B],
) Seq[A, B] {
	var seq Seq[curie.IRI, B]

	switch {
	case pred == nil:
		seq = skiplist.Values(list)
	case pred.Clause == hexer.EQ:
		seq = skiplist.Slice(list, pred.Value, 1)
	case pred.Clause == hexer.PQ:
		_, after := skiplist.Split(list, pred.Value)
		seq = after.TakeWhile(
			func(x curie.IRI) bool {
				return strings.HasPrefix(string(x), string(pred.Value))
			},
		)
	case pred.Clause == hexer.LT:
		seq, _ = skiplist.Split(list, pred.Value)
	case pred.Clause == hexer.GT:
		_, seq = skiplist.Split(list, pred.Value)
	case pred.Clause == hexer.IN:
		seq = skiplist.Range(list, pred.Value, pred.Other)
	}

	return seq.(Seq[A, B])
}

// helper function to query the skiplist where key is xsd.Value
func queryXSD[A, B any](
	pred *hexer.Predicate[xsd.Value],
	list *skiplist.SkipList[xsd.Value, B],
) Seq[A, B] {
	var seq Seq[xsd.Value, B]

	switch {
	case pred == nil:
		seq = skiplist.Values(list)
	case pred.Clause == hexer.EQ:
		seq = skiplist.Slice(list, pred.Value, 1)
	case pred.Clause == hexer.PQ:
		_, after := skiplist.Split(list, pred.Value)
		seq = after.TakeWhile(
			func(x xsd.Value) bool { return xsd.HasPrefix(x, pred.Value) },
		)
	case pred.Clause == hexer.LT:
		before, _ := skiplist.Split(list, pred.Value)
		seq = NewDropWhileType[B](pred.Value.XSDType(), before)
	case pred.Clause == hexer.GT:
		_, after := skiplist.Split(list, pred.Value)
		seq = NewTakeWhileType[B](pred.Value.XSDType(), after)
	case pred.Clause == hexer.IN:
		seq = skiplist.Range(list, pred.Value, pred.Other)
	}

	return seq.(Seq[A, B])
}

// take sequence elements while xsd.Value belongs to same category (type)
type takeWhileType[T any] struct {
	Seq[xsd.Value, T]
	cat curie.IRI
}

func NewTakeWhileType[T any](cat curie.IRI, seq Seq[xsd.Value, T]) Seq[xsd.Value, T] {
	return &takeWhileType[T]{Seq: seq, cat: cat}
}

func (seq *takeWhileType[T]) Next() bool {
	if !seq.Seq.Next() {
		return false
	}

	if key, _ := seq.Seq.Head(); key.XSDType() != seq.cat {
		return false
	}

	return true
}

type dropWhileType[T any] struct {
	Seq[xsd.Value, T]
	cat curie.IRI
}

func NewDropWhileType[T any](cat curie.IRI, seq Seq[xsd.Value, T]) Seq[xsd.Value, T] {
	return &dropWhileType[T]{Seq: seq, cat: cat}
}

func (seq *dropWhileType[T]) Next() bool {
	for {
		if !seq.Seq.Next() {
			return false
		}

		if key, _ := seq.Seq.Head(); key.XSDType() == seq.cat {
			return true
		}
	}
}

// executes query against ⟨s, p, o⟩ data structure
type querySPO hexer.Query

func (q querySPO) L1(list *skiplist.SkipList[s, _po]) Seq[s, _po] {
	return queryIRI[s](q.Pattern.S, list)
}

func (q querySPO) L2(list *skiplist.SkipList[p, __o]) Seq[p, __o] {
	return queryIRI[p](q.Pattern.P, list)
}

func (q querySPO) L3(list *skiplist.SkipList[o, k]) Seq[o, k] {
	return queryXSD[o](q.Pattern.O, list)
}

func (q querySPO) ToSPOCK(s s, p p, o o) hexer.SPOCK {
	return hexer.SPOCK{S: s, P: p, O: o}
}

// executes query against ⟨s, o, p⟩ data structure
type querySOP hexer.Query

func (q querySOP) L1(list *skiplist.SkipList[s, _op]) Seq[s, _op] {
	return queryIRI[s](q.Pattern.S, list)
}

func (q querySOP) L2(list *skiplist.SkipList[o, __p]) Seq[o, __p] {
	return queryXSD[o](q.Pattern.O, list)
}

func (q querySOP) L3(list *skiplist.SkipList[p, k]) Seq[p, k] {
	return queryIRI[p](q.Pattern.P, list)
}

func (q querySOP) ToSPOCK(s s, o o, p p) hexer.SPOCK {
	return hexer.SPOCK{S: s, P: p, O: o}
}

// executes query against ⟨p, s, o⟩ data structure
type queryPSO hexer.Query

func (q queryPSO) L1(list *skiplist.SkipList[p, _so]) Seq[p, _so] {
	return queryIRI[p](q.Pattern.P, list)
}

func (q queryPSO) L2(list *skiplist.SkipList[s, __o]) Seq[s, __o] {
	return queryIRI[s](q.Pattern.S, list)
}

func (q queryPSO) L3(list *skiplist.SkipList[o, k]) Seq[o, k] {
	return queryXSD[o](q.Pattern.O, list)
}

func (q queryPSO) ToSPOCK(p p, s s, o o) hexer.SPOCK {
	return hexer.SPOCK{S: s, P: p, O: o}
}

// executes query against ⟨p, o, s⟩ data structure
type queryPOS hexer.Query

func (q queryPOS) L1(list *skiplist.SkipList[p, _os]) Seq[p, _os] {
	return queryIRI[p](q.Pattern.P, list)
}

func (q queryPOS) L2(list *skiplist.SkipList[o, __p]) Seq[o, __p] {
	return queryXSD[o](q.Pattern.O, list)
}

func (q queryPOS) L3(list *skiplist.SkipList[s, k]) Seq[s, k] {
	return queryIRI[s](q.Pattern.S, list)
}

func (q queryPOS) ToSPOCK(p p, o o, s s) hexer.SPOCK {
	return hexer.SPOCK{S: s, P: p, O: o}
}

// executes query against ⟨o, p, s⟩ data structure
type queryOPS hexer.Query

func (q queryOPS) L1(list *skiplist.SkipList[o, _ps]) Seq[o, _ps] {
	return queryXSD[o](q.Pattern.O, list)
}

func (q queryOPS) L2(list *skiplist.SkipList[p, __s]) Seq[p, __s] {
	return queryIRI[p](q.Pattern.P, list)
}

func (q queryOPS) L3(list *skiplist.SkipList[s, k]) Seq[s, k] {
	return queryIRI[s](q.Pattern.S, list)
}

func (q queryOPS) ToSPOCK(o o, p p, s s) hexer.SPOCK {
	return hexer.SPOCK{S: s, P: p, O: o}
}

// executes query against ⟨o, s, p⟩ data structure
type queryOSP hexer.Query

func (q queryOSP) L1(list *skiplist.SkipList[o, _ps]) Seq[o, _ps] {
	return queryXSD[o](q.Pattern.O, list)
}

func (q queryOSP) L2(list *skiplist.SkipList[s, __p]) Seq[s, __p] {
	return queryIRI[s](q.Pattern.S, list)
}

func (q queryOSP) L3(list *skiplist.SkipList[p, k]) Seq[p, k] {
	return queryIRI[p](q.Pattern.P, list)
}

func (q queryOSP) ToSPOCK(o o, s s, p p) hexer.SPOCK {
	return hexer.SPOCK{S: s, P: p, O: o}
}
