//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/hexagon
//

package jsonld_test

import (
	"encoding/json"
	"testing"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid/v2"
	"github.com/fogfish/hexer"
	"github.com/fogfish/hexer/encoding/jsonld"
	"github.com/fogfish/it/v2"
)

func TestJsonLdUnmarshal(t *testing.T) {
	guid.Clock = guid.NewClockMock()
	luid := curie.IRI("_:5...............")

	Codec := func(t *testing.T, input string) it.SeqOf[hexer.SPOCK] {
		t.Helper()
		bag := jsonld.Bag{}
		err := json.Unmarshal([]byte(input), &bag)
		it.Then(t).Should(it.Nil(err))

		return it.Seq(bag)
	}

	t.Run("OnlyProperty", func(t *testing.T) {
		it.Then(t).Should(
			Codec(t, `{
				"prop": "title"
			}`).Equal(
				hexer.From(luid, "prop", "title"),
			),
		)
	})

	t.Run("PropertyWithID", func(t *testing.T) {
		it.Then(t).Should(
			Codec(t, `{
					"@id": "id",
					"prop": "title"
				}`).Equal(
				hexer.From("id", "prop", "title"),
			),
		)
	})

	// t.Run("PropertyInt", func(t *testing.T) {
	// 	it.Then(t).Should(
	// 		Codec(t, `{
	// 			"prop": 10
	// 		}`).Equal(
	// 			hexer.From(luid, "prop", 10),
	// 		),
	// 	)
	// })

	// t.Run("PropertyFloat", func(t *testing.T) {
	// 	it.Then(t).Should(
	// 		Codec(t, `{
	// 			"prop": 10.0
	// 		}`).Equal(
	// 			hexer.From(luid, "prop", 10.0),
	// 		),
	// 	)
	// })

	// t.Run("PropertyBool", func(t *testing.T) {
	// 	it.Then(t).Should(
	// 		Codec(t, `{
	// 			"prop": true
	// 		}`).Equal(
	// 			hexer.From(luid, "prop", true),
	// 		),
	// 	)
	// })

	t.Run("PropertyArray", func(t *testing.T) {
		it.Then(t).Should(
			Codec(t, `{
				"prop": ["a", "b", "c"]
			}`).Equal(
				hexer.From(luid, "prop", "a"),
				hexer.From(luid, "prop", "b"),
				hexer.From(luid, "prop", "c"),
			),
		)
	})

	// t.Run("PropertyArrayHeterogenous", func(t *testing.T) {
	// 	it.Then(t).Should(
	// 		Codec(t, `{
	// 			"prop": [1, "b", true]
	// 		}`).Equal(
	// 			hexer.From(luid, "prop", 1),
	// 			hexer.From(luid, "prop", "b"),
	// 			hexer.From(luid, "prop", true),
	// 		),
	// 	)
	// })

	t.Run("ArrayOfObjects", func(t *testing.T) {
		it.Then(t).Should(
			Codec(t, `[
				{"@id": "id", "prop": "a"},
				{"@id": "id", "porp": "b"}
			]`).Equal(
				hexer.From("id", "prop", "a"),
				hexer.From("id", "porp", "b"),
			),
		)
	})

	t.Run("Graph", func(t *testing.T) {
		it.Then(t).Should(
			Codec(t, `{
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
			}`).Equal(
				hexer.From("a", "prop", curie.IRI("b")),
				hexer.From("a", "porp", curie.IRI("c")),
				hexer.From("b", "prop", "title"),
				hexer.From("c", "prop", "title"),
			),
		)
	})
}
