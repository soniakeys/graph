// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleBellmanFord() {
	g := graph.WeightedAdjacencyList{
		1: {{2, 10}, {8, 8}},
		2: {{6, 2}},
		3: {{2, 1}, {4, 1}},
		4: {{5, 3}},
		5: {{6, -1}},
		6: {{3, -2}},
		7: {{6, -1}, {2, -4}},
		8: {{7, 1}},
		9: {{5, -10}, {4, 7}},
	}
	b := graph.NewBellmanFord(g)
	if !b.Run(1) {
		return
	}
	fmt.Println("end    path  path")
	fmt.Println("node:  len   dist")
	for n, p := range b.Result.Paths {
		fmt.Printf("%d:       %d   %4.0f\n", n, p.Len, p.Dist)
	}
	// Output:
	// end    path  path
	// node:  len   dist
	// 0:       0   +Inf
	// 1:       1      0
	// 2:       4      5
	// 3:       6      5
	// 4:       7      6
	// 5:       8      9
	// 6:       5      7
	// 7:       3      9
	// 8:       2      8
	// 9:       0   +Inf
}
