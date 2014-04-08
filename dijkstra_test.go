// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed_test

import (
	"fmt"

	"github.com/soniakeys/ed"
)

func ExampleDijkstra_Path() {
	d := ed.NewDijkstra([][]ed.Half{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1.0}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: {},
	})
	path, dist := d.Path(2, 4)
	fmt.Println("Shortest path:", path)
	fmt.Printf("Path distance: %.1f\n", dist)
	// Output:
	// Shortest path: [{2 +Inf} {3 1.1} {4 0.6}]
	// Path distance: 1.7
}

func ExampleDijkstra_AllPaths() {
	g := [][]ed.Half{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1.0}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: nil,
	}
	d := ed.NewDijkstra(g)
	fmt.Println(d.AllPaths(2), "paths found.")
	// column len is from Result, and will be equal to len(path).
	// column dist is from Result, and will be equal to sum.
	fmt.Println("node:  path                       len  dist   sum")
	for nd := range g {
		r := &d.Result[nd]
		path, dist := d.PathTo(nd)
		fmt.Printf("%d:     %-27s %d   %4.1f  %4.1f\n",
			nd, fmt.Sprint(path), r.PathLen, r.PathDist, dist)
	}
	// Output:
	// 4 paths found.
	// node:  path                       len  dist   sum
	// 0:     []                          0   +Inf  +Inf
	// 1:     []                          0   +Inf  +Inf
	// 2:     [{2 +Inf}]                  1    0.0   0.0
	// 3:     [{2 +Inf} {3 1.1}]          2    1.1   1.1
	// 4:     [{2 +Inf} {3 1.1} {4 0.6}]  3    1.7   1.7
	// 5:     [{2 +Inf} {5 0.2}]          2    0.2   0.2
}

func ExampleDijkstra_PathTo() {
	g := [][]ed.Half{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1.0}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: nil,
	}
	d := ed.NewDijkstra(g)
	d.AllPaths(2)
	path, dist := d.PathTo(4)
	fmt.Println("Shortest path:", path)
	fmt.Printf("Path distance: %.1f\n", dist)
	// Output:
	// Shortest path: [{2 +Inf} {3 1.1} {4 0.6}]
	// Path distance: 1.7
}
