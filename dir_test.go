// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleAdjacencyList_Cyclic() {
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
		3: {},
	}
	fmt.Println(g.Cyclic())
	g[3] = []int{2}
	fmt.Println(g.Cyclic())
	// Output:
	// false
	// true
}

func ExampleAdjacencyList_Topological() {
	g := graph.AdjacencyList{
		1: {2},
		3: {1, 2},
		4: {3, 2},
	}
	fmt.Println(g.Topological())
	g[2] = []int{3}
	fmt.Println(g.Topological())
	// Output:
	// [4 3 1 2 0] []
	// [] [3 2 1]
}

func ExampleAdjacencyList_Tarjan() {
	g := graph.AdjacencyList{
		0: {1},
		1: {4, 2, 5},
		2: {3, 6},
		3: {2, 7},
		4: {5, 0},
		5: {6},
		6: {5},
		7: {3, 6},
	}
	scc := g.Tarjan()
	for _, c := range scc {
		fmt.Println(c)
	}
	fmt.Println(len(scc))
	// Output:
	// [6 5]
	// [7 3 2]
	// [4 1 0]
	// 3
}

func ExampleAdjacencyList_Transpose() {
	g := graph.AdjacencyList{
		2: {0, 1},
	}
	t, m := g.Transpose()
	for n, nbs := range t {
		fmt.Printf("%d: %v\n", n, nbs)
	}
	fmt.Println(m)
	// Output:
	// 0: [2]
	// 1: [2]
	// 2: []
	// 2
}

func ExampleAdjacencyList_EulerianCycle() {
	g := graph.AdjacencyList{
		0: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}
	fmt.Println(g.EulerianCycle())
	// Output:
	// [0 1 2 1 2 2 0] <nil>
}
