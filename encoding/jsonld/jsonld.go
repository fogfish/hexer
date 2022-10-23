package jsonld

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fogfish/curie"
	"github.com/fogfish/guid"
	"github.com/fogfish/hexagon"
)

type Node = map[string]any

// Unmarshal JSON-LD into the store
func Unmarshal(data []byte, store *hexagon.Store) error {
	var bag any

	if err := json.Unmarshal(data, &bag); err != nil {
		return err
	}

	switch val := bag.(type) {
	case []any:
		return decodeArray(store, nil, nil, val)
	case Node:
		graph, has := val["@graph"]
		if has {
			switch seq := graph.(type) {
			case []any:
				return decodeArray(store, nil, nil, seq)
			default:
				return fmt.Errorf("json-ld graph codec do not support %T (%v)", val, val)
			}
		}
		return decodeObject(store, nil, nil, val)
	default:
		return fmt.Errorf("json-ld codec do not support %T (%v)", val, val)
	}
}

// Decode JSON-LD into the store
func Decode(reader io.Reader, store *hexagon.Store) error {
	var bag any

	if err := json.NewDecoder(reader).Decode(&bag); err != nil {
		return err
	}

	switch val := bag.(type) {
	case []any:
		return decodeArray(store, nil, nil, val)
	case Node:
		graph, has := val["@graph"]
		if has {
			switch seq := graph.(type) {
			case []any:
				return decodeArray(store, nil, nil, seq)
			default:
				return fmt.Errorf("json-ld graph codec do not support %T (%v)", val, val)
			}
		}
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
			return fmt.Errorf("json-ld array codec do not support %T (%v)", val, val)
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
		return "", false
	}

	id, ok := raw.(string)
	if !ok {
		return "", false
	}

	return curie.IRI(id), true
}

func decodeObjectProperties(store *hexagon.Store, s curie.IRI, obj Node) error {
	for key, val := range obj {
		if key == "@id" {
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
			if err := decodeNodeObject(store, s, p, o); err != nil {
				return err
			}
		case []any:
			if err := decodeNodeArray(store, s, p, o); err != nil {
				return err
			}
		default:
			return fmt.Errorf("json-ld object codec do not support %T (%v)", val, val)
		}
	}

	return nil
}

func decodeNodeObject(store *hexagon.Store, s, p curie.IRI, node Node) error {
	val, has := node["@value"]
	if has {
		hexagon.Put(store, s, p, val)
		return nil
	}

	iri, has := decodeObjectID(node)
	if has {
		hexagon.Put(store, s, p, iri)
		return nil
	}

	return fmt.Errorf("json-ld node object codec do not support %T (%v)", node, node)
}

func decodeNodeArray(store *hexagon.Store, s, p curie.IRI, array []any) error {
	for _, val := range array {
		switch o := val.(type) {
		case float64:
			hexagon.Put(store, s, p, o)
		case string:
			hexagon.Put(store, s, p, o)
		case bool:
			hexagon.Put(store, s, p, o)
		case Node:
			decodeNodeObject(store, s, p, o)
		default:
			return fmt.Errorf("json-ld node array codec do not support %T (%v)", val, val)
		}
	}

	return nil
}

//
//
//
//

// func From(reader io.Reader, store *hexagon.Store) error {
// 	var ld LinkedData

// 	if err := json.NewDecoder(reader).Decode(&ld); err != nil {
// 		return err
// 	}

// 	if ld.Graph != nil && len(ld.Graph) > 0 {
// 		for _, item := range ld.Graph {
// 			if err := decode(store, item); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// func decodeIRI(node Node) (curie.IRI, bool) {
// 	raw, has := node["@id"]
// 	if !has {
// 		return "", false
// 	}

// 	id, ok := raw.(string)
// 	if !ok {
// 		return "", false
// 	}

// 	return curie.IRI(id), true
// }

// func decodeNode(store *hexagon.Store, s curie.IRI, node Node) error {
// 	for key, val := range node {
// 		if key == "@id" {
// 			continue
// 		}
// 		p := curie.IRI(key)

// 		switch o := val.(type) {
// 		case float64:
// 			hexagon.Put(store, s, p, o)
// 		case string:
// 			hexagon.Put(store, s, p, o)
// 		case bool:
// 			hexagon.Put(store, s, p, o)
// 		case Node:
// 			if err := decodeNodeObject(store, s, p, o); err != nil {
// 				return err
// 			}
// 		case []any:
// 			if err := decodeNodeArray(store, s, p, o); err != nil {
// 				return err
// 			}
// 		default:
// 			return fmt.Errorf("json-ld node codec do not support %T (%v)", val, val)
// 		}
// 	}

// 	return nil
// }

// func decode(store *hexagon.Store, item Node) error {
// 	id, has := decodeIRI(item)
// 	if !has {
// 		id = curie.New("_:%s", guid.L.K(guid.Clock))
// 	}

// 	return decodeNode(store, id, item)
// }
