// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math/big"

	"github.com/soniakeys/graph"
)

func ExampleAdjacencyList_BoundsOk() {
	var g graph.AdjacencyList
	ok, _, _ := g.BoundsOk() // zero value adjacency list is valid
	fmt.Println(ok)
	g = graph.AdjacencyList{
		0: {9},
	}
	fmt.Println(g.BoundsOk()) // arc 0 to 9 invalid with only one node
	// Output:
	// true
	// false 0 9
}

func ExampleOneBits() {
	g := make(graph.AdjacencyList, 5)
	var b big.Int
	fmt.Printf("%b\n", graph.OneBits(&b, len(g)))
	// Output:
	// 11111
}
	
func ExampleAdjacencyList_Simple() {
	g := graph.AdjacencyList{
		2: {0, 1},
	}
	fmt.Println(g.Simple())
	g[1] = []int{1} // loop
	fmt.Println(g.Simple())
	g[1] = nil
	g[2] = append(g[2], 0) // parallel arc
	fmt.Println(g.Simple())
	// Output:
	// true -1
	// false 1
	// false 2
}
