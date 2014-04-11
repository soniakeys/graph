package ed_test

import (
	"fmt"
	"strconv"

	"github.com/soniakeys/ed"
)

func ExampleWeightedFromTree_PathTo() {
	g := ed.WeightedAdjacencyList{
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

func ExampleAdjacencyList_Undirected() {
	g := ed.AdjacencyList{
		0: {1, 2},
		2: {0},
	}
	fmt.Println(g.Undirected())
	g[1] = append(g[1], 0)
	fmt.Println(g.Undirected())
	// Output:
	// false 0 1
	// true -1 -1
}

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

func ExampleAdjacencyList_Bipartite() {
	g := ed.AdjacencyList{
		0: {3},
		1: {3},
		2: {3, 4},
		3: {0, 1, 2},
		4: {2},
	}
	b, c1, c2, oc := g.Bipartite(0)
	if b {
		fmt.Println(
			strconv.FormatInt(c1.Int64(), 2),
			strconv.FormatInt(c2.Int64(), 2))
	}
	g[3] = append(g[3], 4)
	g[4] = append(g[4], 3)
	b, c1, c2, oc = g.Bipartite(0)
	if !b {
		fmt.Println(oc)
	}
	// Output:
	// 111 11000
	// [3 4 2]
}
