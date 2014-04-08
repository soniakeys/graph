package ed_test

import (
	"fmt"

	"github.com/soniakeys/ed"
)

func ExampleWeightedFromTree_PathTo() {
	g := [][]ed.Half{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1.0}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: nil,
	}
	d := ed.NewDijkstra(g)
	fmt.Println("From 2")
	d.AllPaths(2)
	path, dist := d.Result.PathTo(3)
	fmt.Printf("To 3: %v %.1f\n", path, dist)
	path, dist = d.Result.PathTo(4)
	fmt.Printf("To 4: %v %.1f\n", path, dist)
	// Output:
	// From 2
	// To 3: [{2 +Inf} {3 1.1}] 1.1
	// To 4: [{2 +Inf} {3 1.1} {4 0.6}] 1.7
}
