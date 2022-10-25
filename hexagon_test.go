package hexagon_test

import (
	"testing"

	"github.com/fogfish/curie"
	"github.com/fogfish/hexagon"
	"github.com/fogfish/it/v2"
)

type spo struct {
	s curie.IRI
	p curie.IRI
	o string
}

var fixture []spo = []spo{
	// 0:7
	{"ex:1", "name", "a"},
	{"ex:1", "name", "b"},
	{"ex:1", "name", "c"},
	{"ex:1", "name", "d"},
	{"ex:1", "prop", "e"},
	{"ex:1", "prop", "f"},
	{"ex:1", "prop", "g"},
	{"ex:1", "prop", "h"},

	// 8:15
	{"ex:2", "name", "a"},
	{"ex:2", "name", "b"},
	{"ex:2", "name", "c"},
	{"ex:2", "name", "d"},
	{"ex:2", "prop", "e"},
	{"ex:2", "prop", "f"},
	{"ex:2", "prop", "g"},
	{"ex:2", "prop", "h"},
}

func toSeq(seq hexagon.Stream) []spo {
	val := []spo{}
	seq.FMap(func(s, p curie.IRI, o any) error {
		val = append(val, spo{s: s, p: p, o: o.(string)})
		return nil
	})
	return val
}

func TestQuery(t *testing.T) {
	store := hexagon.New()

	t.Run("Put", func(t *testing.T) {
		for _, x := range fixture {
			hexagon.Put(store, x.s, x.p, x.o)
		}
	})

	t.Run("∅ ⇒ spo", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				nil, nil, nil,
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(fixture...))
	})

	t.Run("(s) ⇒ po", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				hexagon.IRI.Eq("ex:2"), nil, nil,
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(fixture[8:]...))
	})

	t.Run("(p) ⇒ so", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				nil, hexagon.IRI.Eq("prop"), nil,
			),
		)

		it.Then(t).
			Should(it.Seq(seq[0:4]).Equal(fixture[4:8]...)).
			Should(it.Seq(seq[4:7]).Equal(fixture[12:15]...))
	})

	t.Run("(o) ⇒ sp", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				nil, nil, hexagon.Eq("b"),
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(fixture[1], fixture[9]))
	})

	t.Run("(s)ᴾ ⇒ o", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				hexagon.IRI.Eq("ex:2"), hexagon.IRI.Lt("prop"), nil,
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(fixture[8:12]...))
	})

	t.Run("(s)ᴾ ⇒ o", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				hexagon.IRI.Eq("ex:2"),
				hexagon.IRI.Lt("prop"),
				hexagon.Lt("c"),
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(fixture[8], fixture[9]))
	})

	t.Run("(s)ᴾ ⇒ ∅", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				hexagon.IRI.Eq("ex:2"),
				hexagon.IRI.Lt("prop"),
				hexagon.Lt("a"),
			),
		)

		it.Then(t).
			Should(it.Seq(seq).BeEmpty())
	})

	t.Run("(s)º ⇒ p", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				hexagon.IRI.Eq("ex:2"),
				nil,
				hexagon.Lt("g"),
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(fixture[8:14]...))
	})

	t.Run("(p)º ⇒ s", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				nil,
				hexagon.IRI.Eq("name"),
				hexagon.Lt("g"),
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(
				fixture[0], fixture[8],
				fixture[1], fixture[9],
				fixture[2], fixture[10],
				fixture[3], fixture[11],
			))
	})

	t.Run("(p)º ⇒ s", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				hexagon.IRI.Lt("ex:2"),
				hexagon.IRI.Eq("name"),
				hexagon.Lt("g"),
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(fixture[0:4]...))
	})

	t.Run("(p)ˢ ⇒ o", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				hexagon.IRI.Lt("ex:2"),
				hexagon.IRI.Eq("name"),
				nil,
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(fixture[0:4]...))
	})

	t.Run("(o)ˢ ⇒ p", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				hexagon.IRI.Lt("ex:2"),
				nil,
				hexagon.Eq("c"),
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(fixture[2]))
	})

	t.Run("(o)ˢ ⇒ p", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				hexagon.IRI.Lt("ex:3"),
				hexagon.IRI.Lt("prop"),
				hexagon.Eq("c"),
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(fixture[2], fixture[10]))
	})

	t.Run("(o)ᴾ ⇒ s", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				nil,
				hexagon.IRI.Lt("prop"),
				hexagon.Eq("c"),
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(fixture[2], fixture[10]))
	})

	t.Run("(sp) ⇒ o", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				hexagon.IRI.Eq("ex:2"),
				hexagon.IRI.Eq("prop"),
				nil,
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(fixture[12:]...))
	})

	t.Run("(so) ⇒ p", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				hexagon.IRI.Eq("ex:1"),
				nil,
				hexagon.Eq("e"),
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(fixture[4]))
	})

	t.Run("(po) ⇒ s", func(t *testing.T) {
		seq := toSeq(
			hexagon.Match(store,
				nil,
				hexagon.IRI.Eq("prop"),
				hexagon.Eq("e"),
			),
		)

		it.Then(t).
			Should(it.Seq(seq).Equal(fixture[4], fixture[12]))
	})
}

func BenchmarkXxx(b *testing.B) {
	store := hexagon.New()
	for _, x := range fixture {
		hexagon.Put(store, x.s, x.p, x.o)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		seq := hexagon.Match(store,
			hexagon.IRI.Eq("ex:2"), nil, nil,
		)
		seq.FMap(func(s, p curie.IRI, o any) error { return nil })
	}
}
