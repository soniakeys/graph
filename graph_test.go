// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleAdjacencyList_HasParallelSort_parallelArcs() {
	g := graph.AdjacencyList{
		1: {0, 0},
	}
	// result true 1 0 means parallel arcs from node 1 to node 0
	fmt.Println(g.HasParallelSort())
	// Output:
	// true 1 0
}

func ExampleAdjacencyList_HasParallelSort_noParallelArcs() {
	g := graph.AdjacencyList{
		1: {0},
	}
	// result false -1 -1 means no parallel arc
	fmt.Println(g.HasParallelSort())
	// Output:
	// false -1 -1
}
