// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleLabeledAdjacencyList_NegativeArc() {
	g := graph.LabeledAdjacencyList{
		2: {{To: 0, Label: 0}, {To: 1, Label: 1}},
	}
	arcWeights := []float64{0, .5}
	w := func(label int) float64 { return arcWeights[label] }
	fmt.Println(g.NegativeArc(w))
	g[0] = []graph.Half{{To: 1, Label: len(arcWeights)}}
	arcWeights = append(arcWeights, -2)
	fmt.Println(g.NegativeArc(w))
	// Output:
	// false
	// true
}

func ExampleLabeledAdjacencyList_BoundsOk() {
	g := graph.LabeledAdjacencyList{
		0: {{0, -1}},
	}
	ok, _, _ := g.BoundsOk()
	fmt.Println(ok)
	g = graph.LabeledAdjacencyList{
		0: {{-1, -1}},
	}
	fmt.Println(g.BoundsOk())
	g = graph.LabeledAdjacencyList{
		0: {{9, -1}},
	}
	fmt.Println(g.BoundsOk())
	// Output:
	// true
	// false 0 {-1 -1}
	// false 0 {9 -1}
}

func ExampleLabeledAdjacencyList_IsUndirected() {
	// multigraph, edges with different labels
	//               ----0
	//  (Label: 0)  /   /  (Label: 1)
	//             1----
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 0}, {To: 1, Label: 1}},
		1: {{To: 0, Label: 0}, {To: 0, Label: 1}},
	}
	ud, _, _ := g.IsUndirected()
	fmt.Println(ud)
	// directed graph, arcs with different labels
	//               --->0
	//  (Label: 0)  /   /  (Label: 1)
	//             1<---
	g = graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 1}},
		1: {{To: 0, Label: 0}},
	}
	fmt.Println(g.IsUndirected())
	// Output:
	// true
	// false 0 {1 1}
}

// A directed graph with negative arc weights.
// Arc weights are encoded simply as label numbers.
func ExampleLabeledAdjacencyList_FloydWarshall_negative() {
	g := graph.LabeledAdjacencyList{
		0: {{To: 2, Label: -1}},
		1: {{To: 3, Label: -2}},
		2: {{To: 1, Label: 4}, {To: 3, Label: 3}},
		3: {{To: 0, Label: 2}},
	}
	d := g.FloydWarshall(func(l int) float64 { return float64(l) })
	for _, di := range d {
		fmt.Printf("%2.0f\n", di)
	}
	// Output:
	// [ 0  3 -1  1]
	// [ 0  0 -1 -2]
	// [ 4  4  0  2]
	// [ 2  5  1  0]
}
