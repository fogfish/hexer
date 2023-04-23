package hexer

import (
	"github.com/fogfish/hexer/xsd"
)

type Object interface{ isObject() }

type XSDAnyURI struct{ Value xsd.AnyURI }

func (XSDAnyURI) isObject() {}

type XSDString struct{ Value xsd.String }

func (XSDString) isObject() {}

// type RDFSBag struct{ Value []Object }

// From builds Object from Golang type
func From[T xsd.DataType](value T) Object {
	switch v := any(value).(type) {
	case string:
		return XSDString{Value: v}
	default:
		panic("xxxx")
	}
}
