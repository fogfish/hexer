//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/hexagon
//

package json

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
		return decodeObject(bag, nil, nil, val)
	default:
		return fmt.Errorf("json codec do not support %T (%v)", val, val)
	}
}

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
			return fmt.Errorf("json array codec do not support %T (%v)", val, val)
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
		raw, has = obj["id"]
		if !has {
			return "", false
		}
	}

	id, ok := raw.(string)
	if !ok {
		return "", false
	}

	return curie.IRI(id), true
}

func decodeObjectProperties(bag *Bag, s curie.IRI, obj map[string]any) error {
	for key, val := range obj {
		if key == "@id" || key == "id" {
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
			if err := decodeObject(bag, &s, &p, o); err != nil {
				return err
			}
		case []any:
			if err := decodeArray(bag, &s, &p, o); err != nil {
				return err
			}
		default:
			return fmt.Errorf("json object codec do not support %T (%v)", val, val)
		}
	}

	return nil
}
