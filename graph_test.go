// Copyright 2017 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

// A directed graph with negative arc weights.
// Arc weights are encoded simply as label numbers.
func ExampleDistanceMatrix_FloydWarshall_labeledDirected() {
	//   (-1)   (4)
	//  0---->2---->1
	//  ^     |     |
	//  |(2)  |(3)  |(-2)
	//  |     v     |
	//  ------3<-----
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 2, Label: -1}},
		1: {{To: 3, Label: -2}},
		2: {{To: 1, Label: 4}, {To: 3, Label: 3}},
		3: {{To: 0, Label: 2}},
	}}
	d := g.DistanceMatrix(func(l graph.LI) float64 { return float64(l) })
	d.FloydWarshall()
	for _, di := range d {
		fmt.Printf("%2.0f\n", di)
	}
	// Output:
	// [ 0  3 -1  1]
	// [ 0  0 -1 -2]
	// [ 4  4  0  2]
	// [ 2  5  1  0]
}

func ExampleDistanceMatrix_FloydWarshall_weightedEdgeList() {
	//    (1)   (4)
	//  0-----2-----1
	//  |     |     |
	//  |(2)  |(3)  |(2)
	//  |     |     |
	//  ------3------
	l := graph.WeightedEdgeList{
		Order:      4,
		WeightFunc: func(l graph.LI) float64 { return float64(l) },
		Edges: []graph.LabeledEdge{
			{graph.Edge{0, 2}, 1},
			{graph.Edge{1, 3}, 2},
			{graph.Edge{2, 1}, 4},
			{graph.Edge{2, 3}, 3},
			{graph.Edge{3, 0}, 2},
		},
	}
	d := l.DistanceMatrix()
	d.FloydWarshall()
	for _, di := range d {
		fmt.Printf("%2.0f\n", di)
	}
	// Output:
	// [ 0  4  1  2]
	// [ 4  0  4  2]
	// [ 1  4  0  3]
	// [ 2  2  3  0]
}

func ExampleDistanceMatrix_FloydWarshallPaths() {
	//   (-1)   (4)
	//  0---->2---->1
	//  ^     |     |
	//  |(2)  |(3)  |(-2)
	//  |     v     |
	//  ------3<-----
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 2, Label: -1}},
		1: {{To: 3, Label: -2}},
		2: {{To: 1, Label: 4}, {To: 3, Label: 3}},
		3: {{To: 0, Label: 2}},
	}}
	d := g.DistanceMatrix(func(l graph.LI) float64 { return float64(l) })
	p := d.FloydWarshallPaths()
	fmt.Println("Distances:")
	for _, di := range d {
		fmt.Printf("%2.0f\n", di)
	}
	fmt.Println("Paths:")
	for i := range p {
		for j := range p {
			fmt.Printf("%d->%d:", i, j)
			p.Path(graph.NI(i), graph.NI(j), func(n graph.NI) {
				fmt.Print(" ", n)
			})
			fmt.Println()
		}
	}
	// Output:
	// Distances:
	// [ 0  3 -1  1]
	// [ 0  0 -1 -2]
	// [ 4  4  0  2]
	// [ 2  5  1  0]
	// Paths:
	// 0->0: 0
	// 0->1: 0 2 1
	// 0->2: 0 2
	// 0->3: 0 2 1 3
	// 1->0: 1 3 0
	// 1->1: 1
	// 1->2: 1 3 0 2
	// 1->3: 1 3
	// 2->0: 2 1 3 0
	// 2->1: 2 1
	// 2->2: 2
	// 2->3: 2 1 3
	// 3->0: 3 0
	// 3->1: 3 0 2 1
	// 3->2: 3 0 2
	// 3->3: 3
}

func ExampleDistanceMatrix_FloydWarshallFromLists() {
	//   (-1)   (4)
	//  0---->2---->1
	//  ^     |     |
	//  |(2)  |(3)  |(-2)
	//  |     v     |
	//  ------3<-----
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 2, Label: -1}},
		1: {{To: 3, Label: -2}},
		2: {{To: 1, Label: 4}, {To: 3, Label: 3}},
		3: {{To: 0, Label: 2}},
	}}
	d := g.DistanceMatrix(func(l graph.LI) float64 { return float64(l) })
	l := d.FloydWarshallFromLists()
	fmt.Println("Distances:")
	for _, di := range d {
		fmt.Printf("%2.0f\n", di)
	}
	fmt.Println("Paths:")
	s := make([]graph.NI, len(d))
	for i, li := range l {
		for j := range l {
			p := li.PathTo(graph.NI(j), s)
			if p[0] != graph.NI(i) {
				p = nil
			}
			fmt.Printf("%d->%d: %d\n", i, j, p)
		}
	}
	// Output:
	// Distances:
	// [ 0  3 -1  1]
	// [ 0  0 -1 -2]
	// [ 4  4  0  2]
	// [ 2  5  1  0]
	// Paths:
	// 0->0: [0]
	// 0->1: [0 2 1]
	// 0->2: [0 2]
	// 0->3: [0 2 1 3]
	// 1->0: [1 3 0]
	// 1->1: [1]
	// 1->2: [1 3 0 2]
	// 1->3: [1 3]
	// 2->0: [2 1 3 0]
	// 2->1: [2 1]
	// 2->2: [2]
	// 2->3: [2 1 3]
	// 3->0: [3 0]
	// 3->1: [3 0 2 1]
	// 3->2: [3 0 2]
	// 3->3: [3]
}
