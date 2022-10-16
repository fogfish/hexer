package hexagon

import (
	"testing"

	"github.com/fogfish/curie"
	"github.com/fogfish/it/v2"
)

func TestQuery(t *testing.T) {
	store := New()

	t.Run("∅ ⇒ spo", func(t *testing.T) {
		p := pattern{
			store: store,
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0000)).
			Should(it.Equal(strategy.String(), "∅ ⇒ spo")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(s) ⇒ po", func(t *testing.T) {
		p := pattern{
			store: store,
			s:     Eq(curie.IRI("some")),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0500)).
			Should(it.Equal(strategy.String(), "(s) ⇒ po")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(p) ⇒ so", func(t *testing.T) {
		p := pattern{
			store: store,
			p:     Eq(curie.IRI("some")),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0050)).
			Should(it.Equal(strategy.String(), "(p) ⇒ so")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(o) ⇒ sp", func(t *testing.T) {
		p := pattern{
			store: store,
			o:     Eq[any]("some"),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0005)).
			Should(it.Equal(strategy.String(), "(o) ⇒ sp")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(s)ᴾ ⇒ o", func(t *testing.T) {
		p := pattern{
			store: store,
			s:     Eq(curie.IRI("some")),
			p:     Lt(curie.IRI("some")),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0510)).
			Should(it.Equal(strategy.String(), "(s)ᴾ ⇒ o")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(s)ᴾ ⇒ o", func(t *testing.T) {
		p := pattern{
			store: store,
			s:     Eq(curie.IRI("some")),
			p:     Lt(curie.IRI("some")),
			o:     Lt[any]("some"),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0510)).
			Should(it.Equal(strategy.String(), "(s)ᴾ ⇒ o")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(s)º ⇒ p", func(t *testing.T) {
		p := pattern{
			store: store,
			s:     Eq(curie.IRI("some")),
			o:     Lt[any]("some"),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0501)).
			Should(it.Equal(strategy.String(), "(s)º ⇒ p")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(p)º ⇒ s", func(t *testing.T) {
		p := pattern{
			store: store,
			p:     Eq(curie.IRI("some")),
			o:     Lt[any]("some"),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0051)).
			Should(it.Equal(strategy.String(), "(p)º ⇒ s")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(p)º ⇒ s", func(t *testing.T) {
		p := pattern{
			store: store,
			p:     Eq(curie.IRI("some")),
			o:     Lt[any]("some"),
			s:     Lt(curie.IRI("some")),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0051)).
			Should(it.Equal(strategy.String(), "(p)º ⇒ s")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(p)ˢ ⇒ o", func(t *testing.T) {
		p := pattern{
			store: store,
			p:     Eq(curie.IRI("some")),
			s:     Lt(curie.IRI("some")),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0150)).
			Should(it.Equal(strategy.String(), "(p)ˢ ⇒ o")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(o)ˢ ⇒ p", func(t *testing.T) {
		p := pattern{
			store: store,
			o:     Eq[any]("some"),
			s:     Lt(curie.IRI("some")),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0105)).
			Should(it.Equal(strategy.String(), "(o)ˢ ⇒ p")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(o)ˢ ⇒ p", func(t *testing.T) {
		p := pattern{
			store: store,
			o:     Eq[any]("some"),
			s:     Lt(curie.IRI("some")),
			p:     Lt(curie.IRI("some")),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0105)).
			Should(it.Equal(strategy.String(), "(o)ˢ ⇒ p")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(o)ᴾ ⇒ s", func(t *testing.T) {
		p := pattern{
			store: store,
			o:     Eq[any]("some"),
			p:     Lt(curie.IRI("some")),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0015)).
			Should(it.Equal(strategy.String(), "(o)ᴾ ⇒ s")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(sp) ⇒ o", func(t *testing.T) {
		p := pattern{
			store: store,
			s:     Eq(curie.IRI("some")),
			p:     Eq(curie.IRI("some")),
			o:     Lt[any]("some"),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0550)).
			Should(it.Equal(strategy.String(), "(sp) ⇒ o")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(so) ⇒ p", func(t *testing.T) {
		p := pattern{
			store: store,
			s:     Eq(curie.IRI("some")),
			p:     Lt(curie.IRI("some")),
			o:     Eq[any]("some"),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0505)).
			Should(it.Equal(strategy.String(), "(so) ⇒ p")).
			ShouldNot(it.Nil(iter))
	})

	t.Run("(po) ⇒ s", func(t *testing.T) {
		p := pattern{
			store: store,
			s:     Lt(curie.IRI("some")),
			p:     Eq(curie.IRI("some")),
			o:     Eq[any]("some"),
		}
		strategy, iter := p.eval()

		it.Then(t).
			Should(it.Equal(strategy, 0055)).
			Should(it.Equal(strategy.String(), "(po) ⇒ s")).
			ShouldNot(it.Nil(iter))
	})

}
