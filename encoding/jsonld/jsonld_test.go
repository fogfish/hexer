package jsonld_test

import (
	"testing"

	"github.com/fogfish/curie"
	"github.com/fogfish/hexagon"
	"github.com/fogfish/hexagon/encoding/jsonld"
	"github.com/fogfish/it/v2"
)

func TestJsonLdUnmarshal(t *testing.T) {
	t.Run("OnlyProperty", func(t *testing.T) {
		input := `{
			"prop": "title"
		}`

		val, err := toNode(input, nil, hexagon.IRI("prop"), nil)

		it.Then(t).
			Should(it.Nil(err)).
			Should(it.Map(val).Have("prop", "title"))
	})

	t.Run("PropertyWithID", func(t *testing.T) {
		input := `{
			"@id": "id",
			"prop": "title"
		}`

		val, err := toNode(input, hexagon.IRI("id"), nil, nil)

		it.Then(t).
			Should(it.Nil(err)).
			Should(it.Map(val).Have("prop", "title"))
	})

	t.Run("PropertyInt", func(t *testing.T) {
		input := `{
			"prop": 10
		}`

		val, err := toNode(input, nil, hexagon.IRI("prop"), nil)

		it.Then(t).
			Should(it.Nil(err)).
			Should(it.Map(val).Have("prop", 10.0))
	})

	t.Run("PropertyFloat", func(t *testing.T) {
		input := `{
			"prop": 10.0
		}`

		val, err := toNode(input, nil, hexagon.IRI("prop"), nil)

		it.Then(t).
			Should(it.Nil(err)).
			Should(it.Map(val).Have("prop", 10.0))
	})

	t.Run("PropertyBool", func(t *testing.T) {
		input := `{
			"prop": true
		}`

		val, err := toNode(input, nil, hexagon.IRI("prop"), nil)

		it.Then(t).
			Should(it.Nil(err)).
			Should(it.Map(val).Have("prop", true))
	})

	t.Run("PropertyArray", func(t *testing.T) {
		input := `{
			"prop": ["a", "b", "c"]
		}`

		val, err := toNode(input, nil, hexagon.IRI("prop"), nil)

		it.Then(t).
			Should(it.Nil(err)).
			Should(it.Map(val).Have("prop", []any{"a", "b", "c"}))
	})

	t.Run("PropertyArrayHeterogenous", func(t *testing.T) {
		input := `{
			"prop": [1, "b", true]
		}`

		val, err := toNode(input, nil, hexagon.IRI("prop"), nil)

		it.Then(t).
			Should(it.Nil(err)).
			Should(it.Map(val).Have("prop", []any{true, 1.0, "b"}))
	})

	t.Run("ArrayOfObjects", func(t *testing.T) {
		input := `[
			{"@id": "id", "prop": "a"},
			{"@id": "id", "porp": "b"}
		]`

		val, err := toNode(input, hexagon.IRI("id"), nil, nil)

		it.Then(t).
			Should(it.Nil(err)).
			Should(it.Map(val).Have("prop", "a")).
			Should(it.Map(val).Have("porp", "b"))
	})

	t.Run("Graph", func(t *testing.T) {
		input := `{
			"@graph": [
				{
					"@id": "a",
					"prop": {"@id": "b"},
					"porp": {"@id": "c"}
				},
				{
					"@id": "b",
					"prop": {"@value": "title"}
				},
				{
					"@id": "c",
					"prop": {"@value": "title"}
				}
			]
		}`

		val, err := toNode(input, hexagon.IRI("a"), nil, nil)
		it.Then(t).
			Should(it.Nil(err)).
			Should(it.Map(val).Have("prop", curie.IRI("b"))).
			Should(it.Map(val).Have("porp", curie.IRI("c")))

		val, err = toNode(input, hexagon.IRI("b"), nil, nil)
		it.Then(t).
			Should(it.Nil(err)).
			Should(it.Map(val).Have("prop", "title"))
	})
}

//
//
// Helper
//
//

type moldIRI = *hexagon.Predicate[curie.IRI]
type moldAny = *hexagon.Predicate[any]

func toNode(input string, s moldIRI, p moldIRI, o moldAny) (hexagon.Node, error) {
	node := hexagon.Node{}
	store := hexagon.New()

	if err := jsonld.Unmarshal([]byte(input), store); err != nil {
		return nil, err
	}

	if err := hexagon.Query(store, s, p, o).FMap(node.Append); err != nil {
		return nil, err
	}

	return node, nil
}
