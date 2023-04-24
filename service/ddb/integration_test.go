// go:build it

//
// go test -tags=it
//

package ddb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fogfish/curie"
	"github.com/fogfish/hexer"
	"github.com/fogfish/hexer/service/ddb"
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
	bag.Ref(D, "follows", B)
	bag.Ref(B, "follows", F)
	bag.Ref(F, "follows", G)
	bag.Ref(D, "relates", G)
	bag.Ref(E, "follows", F)

	bag.Add(B, "status", xsd.From("b"))
	bag.Add(D, "status", xsd.From("d"))
	bag.Add(G, "status", xsd.From("g"))

	return bag
}

func setup(bag hexer.Bag) *ddb.Store {
	store, err := ddb.New("ddb:///thingdb-latest")
	if err != nil {
		panic(err)
	}

	t := time.Now()
	_, err = ddb.Add(context.Background(), store, bag)
	if err != nil {
		panic(err)
	}
	fmt.Printf("==> setup %v\n", time.Since(t))

	return store
}

func TestSocialGraph(t *testing.T) {
	rds := setup(datasetSocialGraph())

	Seq := func(t *testing.T, uid string, req hexer.Query) it.SeqOf[hexer.SPOCK] {
		t.Helper()
		bag := hexer.Bag{}
		seq := ddb.Match(context.Background(), rds, req)
		err := seq.FMap(bag.Join)
		it.Then(t).Should(
			it.Nil(err),
			it.Equal(req.String(), uid),
		)

		return it.Seq(bag)
	}

	t.Run("(s) ⇒ po", func(t *testing.T) {
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

	t.Run("(sp) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sp) ⇒ o",
				hexer.NewQuery(hexer.IRI(C), hexer.IRI("follows"), nil),
			).Equal(
				hexer.Link(C, "follows", B),
				hexer.Link(C, "follows", E),
			),
		)
	})

	t.Run("(sᴾ) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾ) ⇒ o",
				hexer.NewQuery(hexer.IRI(C), hexer.Like("fol"), nil),
			).Equal(
				hexer.Link(C, "follows", B),
				hexer.Link(C, "follows", E),
			),
		)
	})

	// t.Run("(spo) ⇒ ∅", func(t *testing.T) {
	// 	it.Then(t).Should(
	// 		Seq(t, "(spo) ⇒ ∅",
	// 			hexer.NewQuery(hexer.IRI(C), hexer.IRI("follows"), hexer.Ref(E)),
	// 		).Equal(
	// 			hexer.Link(C, "follows", E),
	// 		),
	// 	)
	// })

}
