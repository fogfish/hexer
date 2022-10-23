package json

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid"
	"github.com/fogfish/hexagon"
)

type Node = map[string]any

// Unmarshal JSON into the store
func Unmarshal(data []byte, store *hexagon.Store) error {
	var bag any

	if err := json.Unmarshal(data, &bag); err != nil {
		return err
	}

	switch val := bag.(type) {
	case []any:
		return decodeArray(store, nil, nil, val)
	case Node:
		return decodeObject(store, nil, nil, val)
	default:
		return fmt.Errorf("json codec do not support %T (%v)", val, val)
	}
}

// Decode JSON into the store
func Decode(reader io.Reader, store *hexagon.Store) error {
	var bag any

	if err := json.NewDecoder(reader).Decode(&bag); err != nil {
		return err
	}

	switch val := bag.(type) {
	case []any:
		return decodeArray(store, nil, nil, val)
	case Node:
		return decodeObject(store, nil, nil, val)
	default:
		return fmt.Errorf("json codec do not support %T (%v)", val, val)
	}
}

func decodeArray(store *hexagon.Store, s, p *curie.IRI, seq []any) error {
	for _, val := range seq {
		switch o := val.(type) {
		case float64:
			if s != nil && p != nil {
				hexagon.Put(store, *s, *p, o)
			}
		case string:
			if s != nil && p != nil {
				hexagon.Put(store, *s, *p, o)
			}
		case bool:
			if s != nil && p != nil {
				hexagon.Put(store, *s, *p, o)
			}
		case Node:
			decodeObject(store, s, p, o)
		default:
			return fmt.Errorf("json array codec do not support %T (%v)", val, val)
		}
	}
	return nil
}

func decodeObject(store *hexagon.Store, s, p *curie.IRI, obj Node) error {
	id, has := decodeObjectID(obj)
	if !has {
		id = curie.New("_:%s", guid.L.K(guid.Clock))
	}

	if s != nil && p != nil {
		hexagon.Put(store, *s, *p, id)
	}

	return decodeObjectProperties(store, id, obj)
}

func decodeObjectID(obj Node) (curie.IRI, bool) {
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

func decodeObjectProperties(store *hexagon.Store, s curie.IRI, obj Node) error {
	for key, val := range obj {
		if key == "@id" || key == "id" {
			continue
		}
		p := curie.IRI(key)

		switch o := val.(type) {
		case float64:
			hexagon.Put(store, s, p, o)
		case string:
			hexagon.Put(store, s, p, o)
		case bool:
			hexagon.Put(store, s, p, o)
		case Node:
			if err := decodeObject(store, &s, &p, o); err != nil {
				return err
			}
		case []any:
			if err := decodeArray(store, &s, &p, o); err != nil {
				return err
			}
		default:
			return fmt.Errorf("json object codec do not support %T (%v)", val, val)
		}
	}

	return nil
}
