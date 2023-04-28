package ephemeral_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/fogfish/curie"
	"github.com/fogfish/hexer"
	"github.com/fogfish/hexer/service/ephemeral"
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
	N = curie.IRI("n:N")
)

func datasetSocialGraph() hexer.Bag {
	return hexer.Bag{
		hexer.From(A, "follows", B),
		hexer.From(C, "follows", B),
		hexer.From(C, "follows", E),
		hexer.From(C, "relates", D),
		hexer.From(D, "relates", B),
		hexer.From(B, "follows", F),
		hexer.From(F, "follows", G),
		hexer.From(D, "relates", G),
		hexer.From(E, "follows", F),

		hexer.From(B, "status", "b"),
		hexer.From(D, "status", "d"),
		hexer.From(G, "status", "g"),
	}
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

	Seq := func(t *testing.T, uid string, req hexer.Pattern) it.SeqOf[hexer.SPOCK] {
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

	//
	// #2: (s) ⇒ po
	//
	t.Run("#2: (s) ⇒ po", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(s) ⇒ po",
				hexer.Query(hexer.IRI.Equal(C), nil, nil),
			).Equal(
				hexer.From(C, "follows", B),
				hexer.From(C, "follows", E),
				hexer.From(C, "relates", D),
			),
		)
	})

	t.Run("#2: (s) ⇒ po", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(s) ⇒ po",
				hexer.Query(hexer.IRI.Equal(N), nil, nil),
			).Equal(),
		)
	})

	//
	// #3: (sp) ⇒ o
	//
	t.Run("#3: (sp) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sp) ⇒ o",
				hexer.Query(hexer.IRI.Equal(C), hexer.IRI.Equal("follows"), nil),
			).Equal(
				hexer.From(C, "follows", B),
				hexer.From(C, "follows", E),
			),
		)
	})

	t.Run("#3: (sp) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sp) ⇒ o",
				hexer.Query(hexer.IRI.Equal(C), hexer.IRI.Equal("none"), nil),
			).Equal(),
		)
	})

	t.Run("#3: (sp) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sp) ⇒ o",
				hexer.Query(hexer.IRI.Equal(N), hexer.IRI.Equal("follows"), nil),
			).Equal(),
		)
	})

	//
	// #4: (sᴾ) ⇒ o
	//

	t.Run("#4: (sᴾ) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾ) ⇒ o",
				hexer.Query(hexer.IRI.Equal(C), hexer.IRI.HasPrefix("f"), nil),
			).Equal(
				hexer.From(C, "follows", B),
				hexer.From(C, "follows", E),
			),
		)
	})

	t.Run("#4: (sᴾ) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾ) ⇒ o",
				hexer.Query(hexer.IRI.Equal(C), hexer.IRI.HasPrefix("n"), nil),
			).Equal(),
		)
	})

	t.Run("#4: (sᴾ) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sᴾ) ⇒ o",
				hexer.Query(hexer.IRI.Equal(N), hexer.IRI.HasPrefix("f"), nil),
			).Equal(),
		)
	})

	//
	// #5: (so) ⇒ p
	//

	t.Run("#5: (so) ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(so) ⇒ p",
				hexer.Query(hexer.IRI.Equal(D), nil, hexer.Eq(G)),
			).Equal(
				hexer.From(D, "relates", G),
			),
		)
	})

	t.Run("#5: (so) ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(so) ⇒ p",
				hexer.Query(hexer.IRI.Equal(D), nil, hexer.Eq("d")),
			).Equal(
				hexer.From(D, "status", "d"),
			),
		)
	})

	t.Run("#5: (so) ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(so) ⇒ p",
				hexer.Query(hexer.IRI.Equal(D), nil, hexer.Eq(N)),
			).Equal(),
		)
	})

	t.Run("#5: (so) ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(so) ⇒ p",
				hexer.Query(hexer.IRI.Equal(N), nil, hexer.Eq(G)),
			).Equal(),
		)
	})

	//
	// #6: (sº) ⇒ p
	//

	t.Run("#6: (sº) ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sº) ⇒ p",
				hexer.Query(hexer.IRI.Equal(D), nil, hexer.HasPrefix(curie.IRI("s:"))),
			).Equal(
				hexer.From(D, "relates", G),
			),
		)
	})

	t.Run("#6: (s)º ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(s)º ⇒ p",
				hexer.Query(hexer.IRI.Equal(D), nil, hexer.Gt("a")),
			).Equal(
				hexer.From(D, "status", "d"),
			),
		)
	})

	// TODO: test all functions
	// t.Run("#6: (s)º ⇒ p", func(t *testing.T) {
	// 	it.Then(t).Should(
	// 		Seq(t, "(s)º ⇒ p",
	// 			hexer.NewQuery(hexer.IRI.Equal(D), nil, hexer.Lt(curie.IRI("x"))),
	// 		).Equal(
	// 			hexer.From(D, "status", "d"),
	// 		),
	// 	)
	// })

	t.Run("#7: (spo) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(spo) ⇒ ∅",
				hexer.Query(hexer.IRI.Equal(C), hexer.IRI.Equal("follows"), hexer.Eq(E)),
			).Equal(
				hexer.From(C, "follows", E),
			),
		)
	})

	t.Run("#8: (so)ᴾ ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(soᴾ) ⇒ ∅",
				hexer.Query(hexer.IRI.Equal(C), hexer.IRI.HasPrefix("f"), hexer.Eq(E)),
			).Equal(
				hexer.From(C, "follows", E),
			),
		)
	})

	t.Run("#9: (spº) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(spº) ⇒ ∅",
				hexer.Query(hexer.IRI.Equal(C), hexer.IRI.Equal("follows"), hexer.HasPrefix(curie.IRI("u:"))),
			).Equal(
				hexer.From(C, "follows", B),
				hexer.From(C, "follows", E),
			),
		)
	})

	t.Run("#9: (sp)º ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(sp)º ⇒ ∅",
				hexer.Query(hexer.IRI.Equal(G), hexer.IRI.Equal("status"), hexer.Gt("a")),
			).Equal(
				hexer.From(G, "status", "g"),
			),
		)
	})

	// t.Run("#10: (sᴾ)º ⇒ ∅", func(t *testing.T) {
	// 	it.Then(t).Should(
	// 		Seq(t, "(sᴾ)º ⇒ ∅",
	// 			hexer.NewQuery(hexer.IRI(C), hexer.Like("st"), hexer.Gt("a")),
	// 		).Equal(
	// 			hexer.Link(C, "follows", E),
	// 		),
	// 	)
	// })

	t.Run("#11: (p) ⇒ so", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(p) ⇒ so",
				hexer.Query(nil, hexer.IRI.Equal("status"), nil),
			).Equal(
				hexer.From(G, "status", "g"), // s:G < u:B
				hexer.From(B, "status", "b"),
				hexer.From(D, "status", "d"),
			),
		)
	})

	t.Run("#12: (po) ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(po) ⇒ s",
				hexer.Query(nil, hexer.IRI.Equal("follows"), hexer.Eq(B)),
			).Equal(
				hexer.From(C, "follows", B), // s:G < u:A
				hexer.From(A, "follows", B),
			),
		)
	})

	t.Run("#13: (pº) ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(pº) ⇒ s",
				hexer.Query(nil, hexer.IRI.Equal("follows"), hexer.HasPrefix(curie.IRI("s:"))),
			).Equal(
				hexer.From(B, "follows", F),
				hexer.From(E, "follows", F),
				hexer.From(F, "follows", G),
			),
		)
	})

	t.Run("#14: (pˢ) ⇒ o", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(pˢ) ⇒ o",
				hexer.Query(hexer.IRI.HasPrefix("s:"), hexer.IRI.Equal("follows"), nil),
			).Equal(
				hexer.From(C, "follows", B),
				hexer.From(C, "follows", E),
				hexer.From(F, "follows", G),
			),
		)
	})

	// t.Run("#15: (po)ˢ ⇒ ∅", func(t *testing.T) {
	// 	it.Then(t).Should(
	// 		Seq(t, "(po)ˢ ⇒ ∅",
	// 			hexer.NewQuery(hexer.Like("s:"), hexer.IRI("follows"), hexer.Ref(E)),
	// 		).Equal(
	// 			hexer.From(C, "follows", E),
	// 		),
	// 	)
	// })

	// t.Run("#16: (po)ˢ ⇒ ∅", func(t *testing.T) {
	// 	it.Then(t).Should(
	// 		Seq(t, "(po)ˢ ⇒ ∅",
	// 			hexer.NewQuery(hexer.Like("s:"), hexer.IRI("follows"), hexer.Ref("s:")),
	// 		).Equal(
	// 			hexer.From(C, "follows", E),
	// 		),
	// 	)
	// })

	t.Run("#17: (o) ⇒ ps", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(o) ⇒ ps",
				hexer.Query(nil, nil, hexer.Eq(B)),
			).Equal(
				hexer.From(C, "follows", B),
				hexer.From(A, "follows", B),
				hexer.From(D, "relates", B),
			),
		)
	})

	t.Run("#18: (oᴾ) ⇒ s", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(oᴾ) ⇒ s",
				hexer.Query(nil, hexer.IRI.HasPrefix("f"), hexer.Eq(B)),
			).Equal(
				hexer.From(C, "follows", B),
				hexer.From(A, "follows", B),
			),
		)
	})

	t.Run("#19: (oˢ) ⇒ p", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(oˢ) ⇒ p",
				hexer.Query(hexer.IRI.HasPrefix("u:"), nil, hexer.Eq(B)),
			).Equal(
				hexer.From(A, "follows", B),
				hexer.From(D, "relates", B),
			),
		)
	})

	t.Run("#20: (oᴾˢ) ⇒ ∅", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(oᴾˢ) ⇒ ∅",
				hexer.Query(hexer.IRI.HasPrefix("u:"), hexer.IRI.HasPrefix("f"), hexer.Eq(B)),
			).Equal(
				hexer.From(A, "follows", B),
			),
		)
	})

	t.Run("#21: (ˢ) ⇒ po", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ˢ) ⇒ po",
				hexer.Query(hexer.IRI.HasPrefix("s:"), nil, nil),
			).Equal(
				hexer.From(C, "follows", B),
				hexer.From(C, "follows", E),
				hexer.From(C, "relates", D),
				hexer.From(F, "follows", G),
				hexer.From(G, "status", "g"),
			),
		)
	})

	t.Run("#25: (ᴾ) ⇒ so", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(ᴾ) ⇒ so",
				hexer.Query(nil, hexer.IRI.HasPrefix("rel"), nil),
			).Equal(
				hexer.From(C, "relates", D),
				hexer.From(D, "relates", G),
				hexer.From(D, "relates", B),
			),
		)
	})

	t.Run("#27: (º) ⇒ ps", func(t *testing.T) {
		it.Then(t).Should(
			Seq(t, "(º) ⇒ ps",
				hexer.Query(nil, nil, hexer.HasPrefix(curie.IRI("u:"))),
			).Equal(
				hexer.From(C, "follows", B),
				hexer.From(A, "follows", B),
				hexer.From(D, "relates", B),
				hexer.From(C, "relates", D),
				hexer.From(C, "follows", E),
			),
		)
	})

}
