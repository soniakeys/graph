// Copyright 2016 Sonia Keys
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
	cyclic, _, _ := g.Cyclic()
	fmt.Println(cyclic)
	g[3] = []graph.NI{1}
	fmt.Println(g.Cyclic())
	// Output:
	// false
	// true 3 1
}
