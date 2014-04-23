// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleAStar_AStarAPath() {
	a := graph.NewAStar(graph.WeightedAdjacencyList{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: {},
	})
	h4 := []float64{1.9, 2, 1, .6, 0, .9}
	h := func(from int) float64 { return h4[from] }
	p, l := a.AStarAPath(0, 4, h)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Shortest path: [{0 +Inf} {2 0.9} {3 1.1} {4 0.6}]
	// Path length: 2.6
}

func ExampleAStar_AStarMPath() {
	a := graph.NewAStar(graph.WeightedAdjacencyList{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: {},
	})
	h4 := []float64{1.9, 2, 1, .6, 0, .9}
	h := func(from int) float64 { return h4[from] }
	p, l := a.AStarMPath(0, 4, h)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Shortest path: [{0 +Inf} {2 0.9} {3 1.1} {4 0.6}]
	// Path length: 2.6
}

func ExampleHeuristic_Admissable() {
	g := graph.WeightedAdjacencyList{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: {},
	}
	h4 := []float64{1.9, 2, 1, .6, 0, .9}
	var h graph.Heuristic = func(from int) float64 { return h4[from] }
	fmt.Println(h.Admissable(g, 4))
	// Output:
	// true
}

func ExampleHeuristic_Monotonic() {
	g := graph.WeightedAdjacencyList{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: {},
	}
	h4 := []float64{1.9, 2, 1, .6, 0, .9}
	var h graph.Heuristic = func(from int) float64 { return h4[from] }
	fmt.Println(h.Monotonic(g))
	// Output:
	// true
}
