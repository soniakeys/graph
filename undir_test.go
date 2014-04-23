// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed_test

import (
	"fmt"

	"github.com/soniakeys/ed"
)

func ExampleAdjacencyList_ConnectedComponents() {
	g := ed.AdjacencyList{
		0: {3, 4},
		1: {5},
		3: {0, 4},
		4: {0, 3},
		5: {1},
	}
	fmt.Println(g.ConnectedComponents())
	// Output:
	// [0 1 2]
}
