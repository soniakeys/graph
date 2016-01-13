package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExamplePrim_Span() {
	// graph:
	//
	//  (2)     (3)
	//   |\       \
	//   | \       \ 2
	//   |  \       \
	// 4 |   \ 5    (4)
	//   |    \
	//   |     \
	//   |      \
	//  (1)-----(0)
	//       3
	g := &graph.LabeledAdjacencyList{}
	g.AddEdge(graph.LabeledEdge{graph.Edge{0, 1}, 3})
	g.AddEdge(graph.LabeledEdge{graph.Edge{1, 2}, 4})
	g.AddEdge(graph.LabeledEdge{graph.Edge{2, 0}, 5})
	g.AddEdge(graph.LabeledEdge{graph.Edge{3, 4}, 2})
	w := func(arcLabel int) float64 { return float64(arcLabel) }

	// get connected components
	reps, orders := g.ConnectedComponentReps()
	fmt.Println("Connected components:")
	fmt.Println("representative node - order (number of nodes) in component")
	for i, r := range reps {
		fmt.Printf("%d %21d\n", r, orders[i])
	}

	// construct prim object
	p := graph.NewPrim(*g, w)

	// construct spanning tree for each component
	for _, r := range reps {
		ns := p.Span(r)
		fmt.Printf("From node %d, %d nodes spanned.\n", r, ns)
	}

	fmt.Println("Spanning Forest:")
	fmt.Println("Node  From")
	for n, pe := range p.Tree.Paths {
		fmt.Printf("%d %8d\n", n, pe.From)
	}
	// Output:
	// Connected components:
	// representative node - order (number of nodes) in component
	// 0                     3
	// 3                     2
	// From node 0, 3 nodes spanned.
	// From node 3, 2 nodes spanned.
	// Spanning Forest:
	// Node  From
	// 0       -1
	// 1        0
	// 2        1
	// 3       -1
	// 4        3
}

func ExampleWeightedEdgeList_Kruskal() {
	//       (10)
	//     0------4----\
	//     |     /|     \(70)
	// (30)| (40) |(60)  \
	//     |/     |      |
	//     1------2------3
	//       (50)   (20)
	w := func(l int) float64 { return float64(l) }
	l := graph.WeightedEdgeList{5, w, []graph.LabeledEdge{
		{graph.Edge{0, 1}, 30},
		{graph.Edge{0, 4}, 10},
		{graph.Edge{1, 2}, 50},
		{graph.Edge{1, 4}, 40},
		{graph.Edge{2, 3}, 20},
		{graph.Edge{2, 4}, 60},
		{graph.Edge{3, 4}, 70},
	}}
	f, labels, dist := l.Kruskal()
	fmt.Println("node  from  distance  leaf")
	for n, e := range f.Paths {
		fmt.Printf("%d %8d %9.0f %5d\n",
			n, e.From, w(labels[n]), f.Leaves.Bit(n))
	}
	fmt.Println("total distance: ", dist)
	// Output:
	// node  from  distance  leaf
	// 0       -1         0     0
	// 1        0        30     0
	// 2        1        50     0
	// 3        2        20     1
	// 4        0        10     1
	// total distance:  110
}
