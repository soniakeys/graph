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

func ExampleLabeledAdjacencyList_ValidTo() {
	g := graph.LabeledAdjacencyList{
		0: {{0, -1}},
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

func ExampleLabeledFromTree_PathTo() {
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 7}, {To: 2, Label: 9}, {To: 5, Label: 14}},
		1: {{To: 2, Label: 10}, {To: 3, Label: 15}},
		2: {{To: 3, Label: 11}, {To: 5, Label: 2}},
		3: {{To: 4, Label: 6}},
		4: {{To: 5, Label: 9}},
		5: nil,
	}
	w := func(label int) float64 { return float64(label) }
	d := graph.NewDijkstra(g, w)
	fmt.Println("From 2")
	d.AllPaths(2)
	path := d.Tree.PathTo(3)
	dist := d.Dist[3]
	fmt.Printf("To 3: %v %.0f\n", path, dist)
	path = d.Tree.PathTo(4)
	dist = d.Dist[4]
	fmt.Printf("To 4: %v %.0f\n", path, dist)
	// Output:
	// From 2
	// To 3: [{2 -1} {3 11}] 11
	// To 4: [{2 -1} {3 11} {4 6}] 17
}
