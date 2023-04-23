package main

import "fmt"

type X struct{ Value string }

func (x *X) V() { fmt.Printf("%v", x) }

type T struct{ X *X }

func main() {
	x := 05 << 6
	x = x | 07

	fmt.Printf("%b\n", x)
}
