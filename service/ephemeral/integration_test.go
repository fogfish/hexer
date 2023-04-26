package ephemeral_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/fogfish/curie"
	"github.com/fogfish/hexer"
	"github.com/fogfish/hexer/service/ephemeral"
	"github.com/fogfish/hexer/xsd"
	"github.com/fogfish/it/v2"
)

const (
	A = curie.IRI("u:A")
	B = curie.IRI("u:B")
	C = curie.IRI("s:C")
	D = curie.IRI("u:D")
	E = curie.IRI("u:E")
	F = curie.IRI("s:F")
	G = curie.IRI("s:G")
)

func datasetSocialGraph() hexer.Bag {
	bag := hexer.Bag{}
	bag.Ref(A, "follows", B)
	bag.Ref(C, "follows", B)
	bag.Ref(C, "follows", E)
	bag.Ref(C, "relates", D)
	bag.Ref(D, "relates", B)
	bag.Ref(B, "follows", F)
	bag.Ref(F, "follows", G)
	bag.Ref(D, "relates", G)
	bag.Ref(E, "follows", F)

	bag.Add(B, "status", xsd.From("b"))
	bag.Add(D, "status", xsd.From("d"))
	bag.Add(G, "status", xsd.From("g"))

	return bag
}

func setup(bag hexer.Bag) *ephemeral.Store {
	store := ephemeral.New()

	t := time.Now()
	ephemeral.Add(store, bag)

	fmt.Printf("==> setup %v\n", time.Since(t))

	return store
}

func TestSocialGraph(t *testing.T) {
	rds := setup(datasetSocialGraph())

	Seq := func(t *testing.T, uid string, req hexer.Query) it.SeqOf[hexer.SPOCK] {
		t.Helper()
		bag := hexer.Bag{}
		seq := ephemeral.Match(rds, req)
		err := seq.FMap(bag.Join)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(req.String(), uid),
		)

		return it.Seq(bag)
	}

	t.Run("#2: (s) ⇒ po", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(s) ⇒ po",
				hexer.NewQuery(hexer.IRI(C), nil, nil),
			).Equal(
				hexer.Link(C, "follows", B),
				hexer.Link(C, "follows", E),
				hexer.Link(C, "relates", D),
			),
		)
	})

	t.Run("#3: (sp) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sp) ⇒ o",
				hexer.NewQuery(hexer.IRI(C), hexer.IRI("follows"), nil),
			).Equal(
				hexer.Link(C, "follows", B),
				hexer.Link(C, "follows", E),
			),
		)
	})

}
