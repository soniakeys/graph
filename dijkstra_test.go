// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed_test

import (
	"fmt"

	"github.com/soniakeys/ed"
)

func ExampleDijkstra_SingleShortestPath() {
	d := ed.NewDijkstra([][]ed.Half{
		0: {{1, 7}, {2, 9}, {5, 14}},
		1: {{2, 10}, {3, 15}},
		2: {{3, 11}, {5, 2}},
		3: {{4, 6}},
		4: {{5, 9}},
		5: nil,
	})
	p, l := d.SingleShortestPath(2, 4)
	if p == nil {
		fmt.Println("No path")
		return
	}
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Shortest path: [{2 0} {3 11} {4 6}]
	// Path length: 17
}
