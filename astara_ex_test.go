// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed_test

import (
	"fmt"

	"github.com/soniakeys/ed"
)

func ExampleAStarA() {
	a := ed.NewAStar([][]ed.Half{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: {},
	})
	h4 := []float64{1.9, 2, 1, .6, 0, .9}
	h := func(from int) float64 { return h4[from] }
	p, l := a.AStarA(0, 4, h)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Shortest path: [{0 0} {2 0.9} {3 1.1} {4 0.6}]
	// Path length: 2.6
}
