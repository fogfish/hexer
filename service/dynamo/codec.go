package dynamo

import (
	"strings"

	"github.com/fogfish/curie"
	"github.com/fogfish/hexer/xsd"
)

//
// Pair codec
//

func encodeI(a curie.IRI) string {
	return string(a)
}

func encodeII(a, b curie.IRI) string {
	return string(a) + "|" + string(b)
}

func decodeII(val string) (curie.IRI, curie.IRI) {
	seq := strings.SplitN(val, "|", 2)
	return curie.IRI(seq[0]), curie.IRI(seq[1])
}

func encodeIV(a curie.IRI, b xsd.Value) string {
	return string(a) + "|" + encodeValue(b)
}

func decodeIV(val string) (curie.IRI, xsd.Value) {
	seq := strings.SplitN(val, "|", 2)
	return curie.IRI(seq[0]), decodeValue(seq[1])
}

func encodeVI(a xsd.Value, b curie.IRI) string {
	return encodeValue(a) + "|" + string(b)
}

func decodeVI(val string) (xsd.Value, curie.IRI) {
	seq := strings.SplitN(val, "|", 2)
	return decodeValue(seq[0]), curie.IRI(seq[1])
}

//
// Value codec - ᴸᴵᴳ
//

func encodeValue(value xsd.Value) string {
	switch v := value.(type) {
	case xsd.AnyURI:
		return "ᴵ" + string(v)
	case xsd.String:
		return "ᴸ" + string(v)
	default:
		panic("not supported")
	}
}

func decodeValue(value string) xsd.Value {
	switch value[:3] {
	case "ᴵ":
		return xsd.AnyURI(curie.IRI(value[3:]))
	case "ᴸ":
		return xsd.String(value[3:])
	}

	return nil
}
