// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math/big"

	"github.com/soniakeys/graph"
)

func ExampleNextOne() {
	var b big.Int
	b.SetString(""+
		"0000000000000003"+ // bit postions 128, 129
		"0000000000000000"+
		"0000000000000005", 16) // bit positions 2, 0
	for n := graph.NextOne(&b, 0); n >= 0; n = graph.NextOne(&b, n+1) {
		fmt.Println(n)
	}
	// Output:
	// 0
	// 2
	// 128
	// 129
}

func ExampleOneBits() {
	g := make(graph.AdjacencyList, 5)
	var b big.Int
	fmt.Printf("%b\n", graph.OneBits(&b, len(g)))
	// Output:
	// 11111
}

func ExamplePopCount() {
	var b big.Int
	b.SetString(""+
		"0000000000000001"+
		"0000000000000000"+
		"0000000000000005", 16)
	fmt.Println(graph.PopCount(&b))
	// Output:
	// 3
}
