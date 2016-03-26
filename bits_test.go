// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleBits_Iterate() {
	var b graph.Bits
	b.SetBit(0, 1)
	b.SetBit(2, 1)
	b.SetBit(128, 1)
	b.SetBit(129, 1)
	b.Iterate(func(n graph.NI) bool {
		fmt.Println(n)
		return true
	})
	// Output:
	// 0
	// 2
	// 128
	// 129
}

func ExampleBits_NextOne() {
	var b graph.Bits
	b.SetBit(0, 1)
	b.SetBit(2, 1)
	b.SetBit(128, 1)
	b.SetBit(129, 1)
	for n := b.NextOne(0); n >= 0; n = b.NextOne(n + 1) {
		fmt.Println(n)
	}
	// Output:
	// 0
	// 2
	// 128
	// 129
}

func ExampleBits_PopCount() {
	var b graph.Bits
	b.SetBit(0, 1)
	b.SetBit(2, 1)
	b.SetBit(128, 1)
	fmt.Println(b.PopCount())
	// Output:
	// 3
}

func ExampleBits_SetAll() {
	g := make(graph.AdjacencyList, 5)
	var b graph.Bits
	b.SetAll(len(g))
	fmt.Printf("%b\n", b)
	// Output:
	// 11111
}
