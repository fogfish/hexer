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
