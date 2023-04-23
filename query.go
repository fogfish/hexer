package hexer

import (
	"github.com/fogfish/curie"
)

//
// The file define query planner for hex
//

// Strategy is a code that defines combination of indexes and resolution
// strategy to be used for query. The code consists of octal digits for
// each index in the order of <s,p,o>:
//   - 0: no filters or constraints are defined
//   - 1: lookup uses filters defined by predicate so that multiple values are inspected
//   - 4: lookup uses prefix match
//   - 5: lookup uses exact by predicate so that single value is inspected
type Strategy int

const (
	STRATEGY_NONE Strategy = iota
	STRATEGY_SPO
	STRATEGY_SOP
	STRATEGY_PSO
	STRATEGY_POS
	STRATEGY_OPS
	STRATEGY_OSP
)

type Hint int

const (
	HINT_NONE Hint = iota
	HINT_FILTER
	HINT_FILTER_PREFIX
	HINT_MATCH
)

type Query struct {
	Strategy                     Strategy
	HintForS, HintForP, HintForO Hint
	Pattern                      Pattern
}

type Pattern struct {
	S *Predicate[curie.IRI]
	P *Predicate[curie.IRI]
	O *Predicate[Object]
}

func NewQuery(s *Predicate[curie.IRI], p *Predicate[curie.IRI], o *Predicate[Object]) Query {
	q := Query{
		Pattern:  Pattern{S: s, P: p, O: o},
		HintForS: hintFor(s),
		HintForP: hintFor(p),
		HintForO: hintFor(o),
	}
	q.Strategy = strategy(q)

	return q
}

func hintFor[T any](pred *Predicate[T]) Hint {
	switch {
	case pred != nil && pred.Clause != EQ && pred.Clause != PQ:
		return HINT_FILTER
	case pred != nil && pred.Clause == PQ:
		return HINT_FILTER_PREFIX
	case pred != nil && pred.Clause == EQ:
		return HINT_MATCH
	default:
		return HINT_NONE
	}
}

// Estimates execution strategy for pattern
func strategy(q Query) Strategy {
	switch {
	case q.HintForS == HINT_MATCH:
		return strategyForS(q)
	case q.HintForP == HINT_MATCH:
		return strategyForP(q)
	case q.HintForO == HINT_MATCH:
		return strategyForO(q)
	// #1: ___ ⇒ spo
	case q.HintForS == HINT_NONE && q.HintForP == HINT_NONE && q.HintForO == HINT_NONE:
		return STRATEGY_NONE
	default:
		return strategyForX(q)
	}
}

func strategyForS(q Query) Strategy {
	switch {
	// #2: x__ ⇒ spo
	case q.HintForP == HINT_NONE && q.HintForO == HINT_NONE:
		return STRATEGY_SPO
	// #3: xx_ ⇒ spo
	case q.HintForP == HINT_MATCH && q.HintForO == HINT_NONE:
		return STRATEGY_SPO
	// #4: xo_ ⇒ spo
	case (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER) && q.HintForO == HINT_NONE:
		return STRATEGY_SPO
	// #5: x_x ⇒ sop
	case q.HintForP == HINT_NONE && q.HintForO == HINT_MATCH:
		return STRATEGY_SOP
	// #6: x_o ⇒ sop
	case q.HintForP == HINT_NONE && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_SOP
	// #7: xxx ⇒ spo
	case q.HintForP == HINT_MATCH && q.HintForO == HINT_MATCH:
		return STRATEGY_SPO
	// #8: xox ⇒ sop
	case (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER) && q.HintForO == HINT_MATCH:
		return STRATEGY_SOP
	// #9: xxo ⇒ spo
	case q.HintForP == HINT_MATCH && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_SPO
	// #10: xoo ⇒ spo
	case (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER) && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_SPO
	default:
		return STRATEGY_NONE
	}
}

func strategyForP(q Query) Strategy {
	switch {
	// #11: _x_ ⇒ pso
	case q.HintForS == HINT_NONE && q.HintForO == HINT_NONE:
		return STRATEGY_PSO
	// #12: _xx ⇒ pos
	case q.HintForS == HINT_NONE && q.HintForO == HINT_MATCH:
		return STRATEGY_POS
	// #13: _xo ⇒ pos
	case q.HintForS == HINT_NONE && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_POS
	// #14: ox_ ⇒ pso
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && q.HintForO == HINT_NONE:
		return STRATEGY_PSO
	// #15: oxx ⇒ pos
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && q.HintForO == HINT_MATCH:
		return STRATEGY_PSO
	// #16: oxo ⇒ pso
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_PSO
	default:
		return STRATEGY_NONE
	}
}

func strategyForO(q Query) Strategy {
	switch {
	// #17: __x ⇒ ops
	case q.HintForS == HINT_NONE && q.HintForP == HINT_NONE:
		return STRATEGY_OPS
	// #18: _ox ⇒ ops
	case q.HintForS == HINT_NONE && (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER):
		return STRATEGY_OPS
	// #19: o_x ⇒ osp
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && q.HintForP == HINT_NONE:
		return STRATEGY_OSP
	// #20: oox ⇒ ops
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER):
		return STRATEGY_OPS
	default:
		return STRATEGY_NONE
	}
}

func strategyForX(q Query) Strategy {
	switch {
	// #21: o__ ⇒ spo
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && q.HintForP == HINT_NONE && q.HintForO == HINT_NONE:
		return STRATEGY_SPO
	// #22: oo_ ⇒ spo
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER) && q.HintForO == HINT_NONE:
		return STRATEGY_SPO
	// #23: o_o ⇒ sop
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && q.HintForP == HINT_NONE && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_SOP
	// #24: ooo  ⇒ spo
	case (q.HintForS == HINT_FILTER_PREFIX || q.HintForS == HINT_FILTER) && (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER) && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_SPO
	// #25: _o_ ⇒ pso
	case q.HintForS == HINT_NONE && (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER) && q.HintForO == HINT_NONE:
		return STRATEGY_PSO
	// #26: _oo ⇒ pos
	case q.HintForS == HINT_NONE && (q.HintForP == HINT_FILTER_PREFIX || q.HintForP == HINT_FILTER) && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_POS
	// #27: __o ⇒ ops
	case q.HintForS == HINT_NONE && q.HintForP == HINT_NONE && (q.HintForO == HINT_FILTER_PREFIX || q.HintForO == HINT_FILTER):
		return STRATEGY_PSO
	default:
		return STRATEGY_NONE
	}
}

/*

	switch {
	// x, x, x ⇒ spo
	case exact(q.S) && exact(q.P) && exact(q.O):
		return STRATEGY_SPO + 0555
	// x, x, r ⇒ spo
	case exact(q.S) && exact(q.P) && affix(q.O):
		return STRATEGY_SPO + 0554
	// x, x, o ⇒ spo
	case exact(q.S) && exact(q.P) && order(q.O):
		return STRATEGY_SPO + 0551
	// x, x, _ ⇒ spo
	case exact(q.S) && exact(q.P) && q.O == nil:
		return STRATEGY_SPO + 0550
	// x, r, x ⇒ spo
	case exact(q.S) && affix(q.P) && exact(q.O):
		return STRATEGY_SOP + 0545
	// x, o, x ⇒ spo
	case exact(q.S) && order(q.P) && exact(q.O):
		return STRATEGY_SOP + 0515
	// x, _, x ⇒ spo
	case exact(q.S) && q.P == nil && exact(q.O):
		return STRATEGY_SOP + 0515
	}


	switch {
	// r, x, x ⇒ pos
	case affix(q.S) && exact(q.P) && exact(q.O):
		return 0455
	// r, x, r ⇒ pso
	case affix(q.S) && exact(q.P) && affix(q.O):
		return 0454
	// r, x, o ⇒ pso
	case affix(q.S) && exact(q.P) && order(q.O):
		return 0451

	// x, x, _ ⇒ spo
	case exact(q.S) && exact(q.S) && q.O == nil:
		return 0550
	// x, _, x ⇒ sop
	case exact(q.S) && q.P == nil && exact(q.O):
		return 0505
	// _, x, x ⇒ pos
	case q.S == nil && exact(q.P) && exact(q.O):
		return 0055
	// r, x, _ ⇒ pso
	case affix(q.S) && exact(q.S) && q.O == nil:
		return 0450
	// r, _, x ⇒ osp
	case affix(q.S) && q.P == nil && exact(q.O):
		return 0505
	// x, _, r ⇒ sop
	case exact(q.S) && q.P == nil && affix(q.O):
		return 0505
	// x, _, o ⇒ sop
	case exact(q.S) && q.P == nil && affix(q.O):
		return 0505

	// x, o, _ ⇒ spo
	case exact(q.S) && order(q.P) && !exact(q.O):
		return 0510
	// x, _, o ⇒ sop
	case exact(q.s) && !exact(q.p) && order(q.o):
		return 0501
	// _, x, o ⇒ pos
	case !exact(q.s) && exact(q.p) && order(q.o):
		return 0051
	// o, x, _ ⇒ pso
	case order(q.s) && exact(q.p) && !exact(q.o):
		return 0150
	// o, _, x ⇒ osp
	case order(q.s) && !exact(q.p) && exact(q.o):
		return 0105
	// _, o, x ⇒ ops
	case !exact(q.s) && order(q.p) && exact(q.o):
		return 0015

	// x, _, _ ⇒ spo
	case exact(q.s) && !exact(q.p) && !exact(q.o):
		return 0500
	// _, x, _ ⇒ pso
	case !exact(q.s) && exact(q.p) && !exact(q.o):
		return 0050
	// _, _, x ⇒ osp
	case !exact(q.s) && !exact(q.p) && exact(q.o):
		return 0005

	// _, _, _ ⇒ spo
	case !exact(q.s) && !exact(q.p) && !exact(q.o):
		return 0000
	}

	return 0777
}
*/
