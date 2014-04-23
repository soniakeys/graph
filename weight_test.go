// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math"

	"github.com/soniakeys/graph"
)

func ExampleWeightedAdjacencyList_NegativeArc() {
	g := graph.WeightedAdjacencyList{
		2: {{0, 0}, {1, .5}},
	}
	fmt.Println(g.NegativeArc())
	g[0] = []graph.Half{{1, -2}}
	fmt.Println(g.NegativeArc())
	// Output:
	// false
	// true
}

func ExampleWeightedAdjacencyList_ValidTo() {
	g := graph.WeightedAdjacencyList{
		0: {{0, math.NaN()}},
	}
	fmt.Println(g.ValidTo())
	g[0][0].To = -1
	fmt.Println(g.ValidTo())
	g[0][0].To = 1
	fmt.Println(g.ValidTo())
	// Output:
	// true
	// false
	// false
}

func ExampleWeightedFromTree_PathTo() {
	g := graph.WeightedAdjacencyList{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1.0}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: nil,
	}
	d := graph.NewDijkstra(g)
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
