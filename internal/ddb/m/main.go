package main

import (
	"context"
	"fmt"

	"github.com/fogfish/hexer"
	"github.com/fogfish/hexer/internal/ddb"
	"github.com/fogfish/hexer/xsd"
)

func main() {
	store, err := ddb.New("ddb:///thingdb-latest")
	if err != nil {
		panic(err)
	}

	bag := hexer.Bag{}
	bag.Ref("A", "follows", "B")
	bag.Ref("C", "follows", "B")
	bag.Ref("C", "follows", "D")
	bag.Ref("D", "follows", "B")
	bag.Ref("B", "follows", "F")
	bag.Ref("F", "follows", "G")
	bag.Ref("D", "follows", "G")
	bag.Ref("E", "follows", "F")

	bag.Add("B", "status", xsd.From("b"))
	bag.Add("D", "status", xsd.From("d"))
	bag.Add("G", "status", xsd.From("g"))

	if _, err := ddb.Add(context.Background(), store, bag); err != nil {
		panic(err)
	}

	q := hexer.NewQuery(hexer.IRI("C"), nil, nil)
	fmt.Printf("==> %s\n", q)

	s := ddb.Match(context.Background(), store, q)
	for s.Next() {
		fmt.Printf("==> %+v\n", s.Head())
	}
}
