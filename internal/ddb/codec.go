package ddb

import (
	"strings"

	"github.com/fogfish/curie"
	"github.com/fogfish/hexer"
)

// func encode[T hexer.DataType](spock hexer.SPOCK[T]) string {
// 	switch v := any(spock.O).(type) {
// 	case string:
// 		return "s" + v
// 	default:
// 		return "g"
// 	}
// }

//
// Pair codec
//

func encodeII(a, b curie.IRI) string {
	return string(a) + "|" + string(b)
}

func decodeII(val string) (curie.IRI, curie.IRI) {
	seq := strings.SplitN(val, "|", 2)
	return curie.IRI(seq[0]), curie.IRI(seq[1])
}

func encodeIV(a curie.IRI, b hexer.XSDValue) string {
	return string(a) + "|" + encodeValue(b)
}

func decodeIV(val string) (curie.IRI, hexer.XSDValue) {
	seq := strings.SplitN(val, "|", 2)
	return curie.IRI(seq[0]), decodeValue(seq[1])
}

func encodeVI(a hexer.XSDValue, b curie.IRI) string {
	return encodeValue(a) + "|" + string(b)
}

func decodeVI(val string) (hexer.XSDValue, curie.IRI) {
	seq := strings.SplitN(val, "|", 2)
	return decodeValue(seq[0]), curie.IRI(seq[1])
}

//
// Value codec
//

func encodeValue(value hexer.XSDValue) string {
	switch v := value.(type) {
	case hexer.XSDAnyURI:
		return string(v.Value)
	case hexer.XSDString:
		return v.Value
	default:
		panic("not supported")
	}
}

func decodeValue(value string) hexer.XSDValue {
	return hexer.XSDString{Value: value}
}
