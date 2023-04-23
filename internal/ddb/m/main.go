package main

import (
	"context"
	"fmt"

	"github.com/fogfish/hexer"
	"github.com/fogfish/hexer/internal/ddb"
)

func main() {
	store, err := ddb.New("ddb:///thingdb-latest")
	if err != nil {
		panic(err)
	}

	seq := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	for _, x := range seq {
		err = ddb.Put(context.Background(), store,
			hexer.SPOCK{
				S: "abc",
				P: "name",
				O: hexer.XSDString{Value: x},
			},
		)
	}

	// q := hexer.NewQuery(hexer.IRI("abc"), hexer.Prefix("na"), nil)
	q := hexer.NewQuery(hexer.Prefix("a"), nil, hexer.Eq("5"))

	s := ddb.Match(context.Background(), store, q)

	for s.Next() {
		fmt.Printf("==> %+v\n", s.Head())
	}

	// store.Iterate()

	// iter := ddb.IT(store)

	// for iter.Next() {
	// 	iter.Head()
	// }

	// spo := ddb.Get(context.Background(), store,
	// 	hexer.SPOCK[string]{
	// 		S: "a",
	// 		P: "name",
	// 	},
	// )
	// fmt.Printf("==> %+v\n", spo)

	// buf := new(bytes.Buffer)
	// enc := gob.NewEncoder(buf)

	// enc.Encode("x:abc")
	// fmt.Printf("%s\n", buf.Bytes())

	// buf.Reset()

	// enc.Encode("x:abe")
	// fmt.Printf("%s\n", buf.Bytes())
}
