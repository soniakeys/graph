// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed_test

import (
	"fmt"
	"testing"

	"github.com/soniakeys/ed"
)

func ExampleDijkstra_Path() {
	d := ed.NewDijkstra([][]ed.Half{
		1: {{2, 7}, {3, 9}, {6, 11}},
		2: {{3, 10}, {4, 15}},
		3: {{4, 11}, {6, 2}},
		4: {{5, 7}},
		6: {{5, 9}},
	})
	path, dist := d.Path(1, 5)
	fmt.Println("Shortest path:", path)
	fmt.Println("Path distance:", dist)
	// Output:
	// Shortest path: [{1 +Inf} {6 11} {5 9}]
	// Path distance: 20
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
		r := &d.Result.Paths[nd]
		path, dist := d.Result.PathTo(nd)
		fmt.Printf("%d:     %-27s %d   %4.1f  %4.1f\n",
			nd, fmt.Sprint(path), r.Len, r.Dist, dist)
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

func BenchmarkDijkstra100(b *testing.B) {
	// 100 nodes, 200 edges
	tc := r100
	d := ed.NewDijkstra(tc.g)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}

func BenchmarkDijkstra1e3(b *testing.B) {
	// 1000 nodes, 3000 edges
	once.Do(bigger)
	tc := r1k
	d := ed.NewDijkstra(tc.g)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}

func BenchmarkDijkstra1e4(b *testing.B) {
	// 10k nodes, 50k edges
	once.Do(bigger)
	tc := r10k
	d := ed.NewDijkstra(tc.g)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}

func BenchmarkDijkstra1e5(b *testing.B) {
	// 100k nodes, 1m edges
	once.Do(bigger)
	tc := r100k
	d := ed.NewDijkstra(tc.g)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}
