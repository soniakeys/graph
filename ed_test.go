// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed_test

import (
	"fmt"

	"github.com/soniakeys/ed"
)

func ExampleGraph_ShortestPath() {
	g := ed.New(6)
	g.SetArcs(0, []ed.Half{{1, 7}, {2, 9}, {5, 14}})
	g.SetArcs(1, []ed.Half{{2, 10}, {3, 15}})
	g.SetArcs(2, []ed.Half{{3, 11}, {5, 2}})
	g.SetArcs(3, []ed.Half{{4, 6}})
	g.SetArcs(4, []ed.Half{{5, 9}})
	p, l := g.ShortestPath(2, 4)
	if p == nil {
		fmt.Println("No path from start node to end node")
		return
	}
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Shortest path: [{2 0} {3 11} {4 6}]
	// Path length: 17
}
