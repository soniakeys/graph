// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"strconv"

	"github.com/soniakeys/graph"
)

func ExampleAdjacencyList_Ok() {
	var g graph.AdjacencyList
	fmt.Println(g.Ok()) // zero value adjacency list is valid
	g = graph.AdjacencyList{
		0: {1},
	}
	fmt.Println(g.Ok()) // arc 0 to 1 invalid with only one node
	// Output:
	// true
	// false
}

func ExampleAdjacencyList_Undirected() {
	g := graph.AdjacencyList{
		0: {1, 2},
		2: {0},
	}
	fmt.Println(g.Undirected())
	g[1] = append(g[1], 0)
	fmt.Println(g.Undirected())
	// Output:
	// false 0 1
	// true -1 -1
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

func ExampleAdjacencyList_Bipartite() {
	g := graph.AdjacencyList{
		0: {3},
		1: {3},
		2: {3, 4},
		3: {0, 1, 2},
		4: {2},
	}
	b, c1, c2, oc := g.Bipartite(0)
	if b {
		fmt.Println(
			strconv.FormatInt(c1.Int64(), 2),
			strconv.FormatInt(c2.Int64(), 2))
	}
	g[3] = append(g[3], 4)
	g[4] = append(g[4], 3)
	b, c1, c2, oc = g.Bipartite(0)
	if !b {
		fmt.Println(oc)
	}
	// Output:
	// 111 11000
	// [3 4 2]
}
