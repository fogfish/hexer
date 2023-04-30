//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/hexagon
//

package jsonld

import (
	"encoding/json"
	"fmt"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid/v2"
	"github.com/fogfish/hexer"
)

type Bag hexer.Bag

func (bag *Bag) UnmarshalJSON(b []byte) error {
	var raw any

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	switch val := raw.(type) {
	case []any:
		return decodeArray(bag, nil, nil, val)
	case map[string]any:
		graph, has := val["@graph"]
		if has {
			switch seq := graph.(type) {
			case []any:
				return decodeArray(bag, nil, nil, seq)
			default:
				return fmt.Errorf("json-ld graph codec do not support %T (%v)", val, val)
			}
		}
		return decodeObject(bag, nil, nil, val)
	default:
		return fmt.Errorf("json-ld codec do not support %T (%v)", val, val)
	}
}

// type Node = map[string]any

// // Unmarshal JSON-LD into the store
// func Unmarshal(data []byte, store *hexagon.Store) error {
// 	var bag any

// 	if err := json.Unmarshal(data, &bag); err != nil {
// 		return err
// 	}

// 	switch val := bag.(type) {
// 	case []any:
// 		return decodeArray(store, nil, nil, val)
// 	case Node:
// 		graph, has := val["@graph"]
// 		if has {
// 			switch seq := graph.(type) {
// 			case []any:
// 				return decodeArray(store, nil, nil, seq)
// 			default:
// 				return fmt.Errorf("json-ld graph codec do not support %T (%v)", val, val)
// 			}
// 		}
// 		return decodeObject(store, nil, nil, val)
// 	default:
// 		return fmt.Errorf("json-ld codec do not support %T (%v)", val, val)
// 	}
// }

func decodeArray(bag *Bag, s, p *curie.IRI, seq []any) error {
	for _, val := range seq {
		switch o := val.(type) {
		case float64:
			if s != nil && p != nil {
				*bag = append(*bag, hexer.From(*s, *p, o))
			}
		case string:
			if s != nil && p != nil {
				*bag = append(*bag, hexer.From(*s, *p, o))
			}
		case bool:
			if s != nil && p != nil {
				*bag = append(*bag, hexer.From(*s, *p, o))
			}
		case map[string]any:
			decodeObject(bag, s, p, o)
		default:
			return fmt.Errorf("json-ld array codec do not support %T (%v)", val, val)
		}
	}
	return nil
}

func decodeObject(bag *Bag, s, p *curie.IRI, obj map[string]any) error {
	id, has := decodeObjectID(obj)
	if !has {
		id = curie.New("_:%s", guid.L(guid.Clock))
	}

	if s != nil && p != nil {
		*bag = append(*bag, hexer.From(*s, *p, id))
	}

	return decodeObjectProperties(bag, id, obj)
}

func decodeObjectID(obj map[string]any) (curie.IRI, bool) {
	raw, has := obj["@id"]
	if !has {
		return "", false
	}

	id, ok := raw.(string)
	if !ok {
		return "", false
	}

	return curie.IRI(id), true
}

func decodeObjectProperties(bag *Bag, s curie.IRI, obj map[string]any) error {
	for key, val := range obj {
		if key == "@id" {
			continue
		}
		p := curie.IRI(key)

		switch o := val.(type) {
		case float64:
			*bag = append(*bag, hexer.From(s, p, o))
		case string:
			*bag = append(*bag, hexer.From(s, p, o))
		case bool:
			*bag = append(*bag, hexer.From(s, p, o))
		case map[string]any:
			if err := decodeNodeObject(bag, s, p, o); err != nil {
				return err
			}
		case []any:
			if err := decodeNodeArray(bag, s, p, o); err != nil {
				return err
			}
		default:
			return fmt.Errorf("json-ld object codec do not support %T (%v)", val, val)
		}
	}

	return nil
}

func decodeNodeObject(bag *Bag, s, p curie.IRI, node map[string]any) error {
	val, has := node["@value"]
	if has {
		return decodeValue(bag, s, p, val)
	}

	iri, has := decodeObjectID(node)
	if has {
		*bag = append(*bag, hexer.From(s, p, iri))
		return nil
	}

	return fmt.Errorf("json-ld node object codec do not support %T (%v)", node, node)
}

func decodeNodeArray(bag *Bag, s, p curie.IRI, array []any) error {
	for _, val := range array {
		switch o := val.(type) {
		case float64:
			*bag = append(*bag, hexer.From(s, p, o))
		case string:
			*bag = append(*bag, hexer.From(s, p, o))
		case bool:
			*bag = append(*bag, hexer.From(s, p, o))
		case map[string]any:
			decodeNodeObject(bag, s, p, o)
		default:
			return fmt.Errorf("json-ld node array codec do not support %T (%v)", val, val)
		}
	}

	return nil
}

func decodeValue(bag *Bag, s, p curie.IRI, val any) error {
	switch o := val.(type) {
	case float64:
		*bag = append(*bag, hexer.From(s, p, o))
	case string:
		*bag = append(*bag, hexer.From(s, p, o))
	case bool:
		*bag = append(*bag, hexer.From(s, p, o))
	default:
		return fmt.Errorf("json-ld value codec do not support %T (%v)", val, val)
	}

	return nil
}
